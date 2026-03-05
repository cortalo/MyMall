package com.longhe.learn.oomall.payment.controller;

import com.longhe.learn.javaee.core.model.ReturnObject;
import com.longhe.learn.oomall.payment.dao.bo.Account;
import com.longhe.learn.oomall.payment.service.ChannelService;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequiredArgsConstructor
@Slf4j
public class CustomerController {

    private final ChannelService channelService;

    /**
     *获得商户有效的收款账号
     */
    @GetMapping("/accounts")
    public ReturnObject getChannels(@RequestParam Long shopId,
                                    @RequestParam(required = false,defaultValue = "1") Integer page,
                                    @RequestParam(required = false,defaultValue = "10") Integer pageSize) {
        List<Account> channels = this.channelService.retrieveValidAccount(shopId, page, pageSize);
        return null;
    }

}
