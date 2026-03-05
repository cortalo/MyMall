package com.longhe.learn.oomall.payment.mapper.provider.wepay;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.NoArgsConstructor;

/**
 * 取消订单请求参数
 * */
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class CancelOrderParam {
    /*必选*/
    /**
     * 直连商户
     * payTrans.account.channel.SpMchid
     * */
    private String sp_mchid;

    /**
     * 子商户的商户号，
     */
    private String sub_mchid;
}
