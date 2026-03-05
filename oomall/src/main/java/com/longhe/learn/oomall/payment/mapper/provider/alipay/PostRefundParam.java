package com.longhe.learn.oomall.payment.mapper.provider.alipay;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.NoArgsConstructor;

import java.util.List;

/**
 * 发起退款(退分账)参数
 * @author Wenbo Li
 * */
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class PostRefundParam {
    /*必选*/
    /**
     * 退款金额
     * */
    private Double refund_amount;

    public void setRefund_amount(Long amount){
        this.refund_amount = amount / 100.0;
    }

    /*二选一*/
    /**
     * 支付交易号
     * refundTrans.payTrans.trans_no
     * */
    private String trade_no;

    /*可选*/
    /**
     * 退款请求号。 标识一次退款请求，需要保证在交易号下唯一
     * refundTrans.outNo
     * */
    private String out_request_no;

    /*退分账*/
    /**
     * 退分账明细
     * */
    private List<RoyaltyDetailInfoPojo> refund_royalty_parameters;

}

