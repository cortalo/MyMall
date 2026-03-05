package com.longhe.learn.oomall.payment.mapper.po;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;
import lombok.ToString;

import java.time.LocalDateTime;

@Entity
@Table(name = "payment_channel")
@Data
@AllArgsConstructor
@NoArgsConstructor
@ToString
public class ChannelPo {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private String spAppid;
    private String name;
    private String spMchid;
    private LocalDateTime beginTime;
    private LocalDateTime endTime;
    private Byte status;
    private String beanName;
    private Long creatorId;
    private String creatorName;
    private Long modifierId;
    private String modifierName;
    private LocalDateTime gmtCreate;
    private LocalDateTime gmtModified;
}
