//School of Informatics Xiamen University, GPL-3.0 license
package com.longhe.learn.oomall.payment.mapper.provider.alipay;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.NoArgsConstructor;

/**
 * 取消订单信息参数
 * @author Wenbo Li
 * */

@NoArgsConstructor
@AllArgsConstructor
@Builder
public class CancelOrderParam {
    /*二选一*/
    /**
     * 支付宝交易号
     * payTrans.trans_no
     * */
    private String trade_no;

}
