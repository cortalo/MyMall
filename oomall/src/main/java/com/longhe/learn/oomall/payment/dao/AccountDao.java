package com.longhe.learn.oomall.payment.dao;

import com.longhe.learn.javaee.core.exception.BusinessException;
import com.longhe.learn.javaee.core.mapper.RedisUtil;
import com.longhe.learn.javaee.core.model.ReturnNo;
import com.longhe.learn.javaee.core.util.CloneFactory;
import com.longhe.learn.oomall.payment.dao.bo.Account;
import com.longhe.learn.oomall.payment.mapper.AccountPoMapper;
import com.longhe.learn.oomall.payment.mapper.po.AccountPo;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Objects;
import java.util.Optional;

import static com.longhe.learn.javaee.core.model.Constants.PLATFORM;

@Repository
@Slf4j
@RequiredArgsConstructor
public class AccountDao {

    private static final String ACCOUNTKEY = "A%d";
    private static final String SHOPACCOUNTKEY = "SA%d";

    @Value("${oomall.payment.account.timeout}")
    private long timeout;

    private final RedisUtil redisUtil;
    private final AccountPoMapper accountPoMapper;
    private final ChannelDao channelDao;

    public List<Account> retrieveByShopId(Long shopId) {
        String key = String.format(SHOPACCOUNTKEY, shopId);
        List<Long> accountIds = (List<Long>) redisUtil.get(key);
        if (Objects.nonNull(accountIds)) {
            return accountIds.stream().map(id -> this.findById(shopId, id)).toList();
        } else {
            log.debug("retrieveByShopId: shopId = {}", shopId);
            return this.accountPoMapper.findByShopId(shopId, Pageable.unpaged()).stream()
                    .map(po -> this.build(po, Optional.empty()))
                    .toList();
        }
    }

    private Account findById(Long shopId, Long id) {
        Account ret = null;
        String key = String.format(ACCOUNTKEY, id);
        ret = (Account) redisUtil.get(key);
        if (!Objects.isNull(ret)) {
            if (!shopId.equals(ret.getShopId()) && !PLATFORM.equals(shopId)) {
                throw new BusinessException(ReturnNo.RESOURCE_ID_OUTSCOPE, String.format(ReturnNo.RESOURCE_ID_OUTSCOPE.getMessage(), "支付渠道", id, shopId));
            }
            this.build(ret);
        } else {
            Optional<AccountPo> optPo = accountPoMapper.findById(id);
            if (optPo.isEmpty()) {
                throw new BusinessException(ReturnNo.RESOURCE_ID_NOTEXIST, String.format(ReturnNo.RESOURCE_ID_NOTEXIST.getMessage(), "商铺支付渠道", id));
            }
            if (!shopId.equals(optPo.get().getShopId()) && !PLATFORM.equals(shopId)) {
                throw new BusinessException(ReturnNo.RESOURCE_ID_OUTSCOPE, String.format(ReturnNo.RESOURCE_ID_OUTSCOPE.getMessage(), "支付渠道", id, shopId));
            }
            ret = this.build(optPo.get(), Optional.of(key));
        }
        log.debug("findObjById: id = " + id + " ret = " + ret);
        return ret;
    }

    private Account build(AccountPo po, Optional<String> redisKey) {
        Account ret = CloneFactory.copy(new Account(), po);
        redisKey.ifPresent(key -> redisUtil.set(key, ret, timeout));
        this.build(ret);
        return ret;
    }

    private Account build(Account bo) {
        bo.setChannelDao(this.channelDao);
        return bo;
    }
}
