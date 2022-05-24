package main

import (
	"errors"
	"github.com/ssbcV2/contract"
	"time"
)

var Beneficiary string      // 众筹受益人
var FundingEnd time.Time    // 结束时间
var Ended bool              // 众筹结束标记
var Gives = map[string]int{} // 所有参与众筹者的出价
var Sum int //当前合约一共收到了多少筹款
var Target int //合约预定需要收到多少筹款才算众筹成功
var ExceedContinue bool //募集到指定资金后是否继续接受新的资金

func init() {
	Beneficiary = contract.Caller()              // 受益人默认为发布合约的人
	FundingEnd = time.Now().Add(time.Minute * 2) // 合约在发布两分钟后停止出价
	Sum = 0 //初始状态下合约中金额为空
	Target = 1000 //由受益人确定需要筹集多少资金，这里默认是0
	exceedContinue = true //由受益人确定众筹到预期金额后是否接收新的资金，这里默认是true（反之为false）
}

func Give(args map[string]interface{}) (interface{}, error) {
	if FundingEnd.Before(time.Now()) {
		contract.Transfer(contract.Caller(), contract.Value()) // 退回转账
		contract.Info("众筹已结束")
		return nil, errors.New("众筹已结束")
	}
	
	if (Sum>=Target && exceedContinue==false){
		contract.Transfer(contract.Caller(), contract.Value()) // 退回转账
		contract.Info("已众筹到预期金额，不再接受新的资金，感谢您的支持！")
		return nil, errors.New("已众筹到预期金额，不再接受新的资金，感谢您的支持！")
	}

	Gives[contract.Caller()] += contract.Value()
	Sum +=  contract.Value()
	return nil, nil
}

func End(args map[string]interface{}) (interface{}, error) {
	contract.Transfer(contract.Caller(), contract.Value()) // FundingEnd方法不接受转账，退回
	if FundingEnd.After(time.Now()) {
		contract.Info("众筹还未结束")
		return nil, errors.New("众筹还未结束")
	}

	if Ended {
		contract.Info("重复调用FundingEnd")
		return nil, errors.New("重复调用FundingEnd")
	}
	Ended = true


	if(Sum >= Target){
		//这种情况下众筹成功，金额从合约账户转账至受益人的账户
		contract.Transfer(Beneficiary, Sum);
	} else {
		for giver, amount := range Gives {
			contract.Transfer(giver, amount) // 众筹失败，资金从众筹账户退回到每个参与众筹者的账户
		}
	}


	return nil, nil
}

// 回退函数，当没有方法匹配时执行此方法
func Fallback(args map[string]interface{}) (interface{}, error) {
	contract.Transfer(contract.Caller(), contract.Value()) // 将转账退回
	return nil, nil
}
