package com.longhe.learn.oomall.payment.service;

import com.longhe.learn.javaee.core.mapper.RedisUtil;
import com.longhe.learn.oomall.payment.dao.AccountDao;
import com.longhe.learn.oomall.payment.dao.bo.Account;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
@Transactional
@RequiredArgsConstructor
@Slf4j
public class ChannelService {

    private final RedisUtil redisUtil;
    private final AccountDao accountDao;

    public List<Account> retrieveValidAccount(Long shopId, Integer page, Integer pageSize) {
        log.debug("retrieveValidAccount: shopId = {}", shopId);
        List<Account> accounts = this.accountDao.retrieveByShopId(shopId);
        log.debug("retrieveValidAccount: accounts = {}", accounts);
        return accounts.stream().filter(account -> account.getStatus().equals(Account.VALID))
                .skip(pageSize * (page - 1)).limit(page * pageSize)
                .toList();
    }
}
