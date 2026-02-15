package com.longhe.learn.mymall.payment.dao.bo;

import com.longhe.learn.mymall.core.model.OOMallObject;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import lombok.ToString;

import java.io.Serializable;

/**
 * 业务
 */
@Data
@ToString(callSuper = true)
@EqualsAndHashCode(callSuper = true)
@NoArgsConstructor
public class Business extends OOMallObject implements Serializable {
    /**
     * 业务名称
     */
    private String name;
}