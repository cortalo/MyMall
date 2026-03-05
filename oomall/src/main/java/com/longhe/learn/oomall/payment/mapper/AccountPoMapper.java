package com.longhe.learn.oomall.payment.mapper;

import com.longhe.learn.oomall.payment.mapper.po.AccountPo;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface AccountPoMapper extends JpaRepository<AccountPo, Long> {
    List<AccountPo> findByShopId(Long shopId, Pageable pageable);
}
