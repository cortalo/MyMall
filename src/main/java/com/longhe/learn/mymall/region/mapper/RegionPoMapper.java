package com.longhe.learn.mymall.region.mapper;

import com.longhe.learn.mymall.region.mapper.po.RegionPo;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface RegionPoMapper extends JpaRepository<RegionPo, Long> {
    List<RegionPo> findByPidEqualsAndStatusEquals(Long pid, Byte status, Pageable pageable);

    List<RegionPo> findByPid(Long pid, Pageable pageable);

    Optional<RegionPo> findByName(String name);
}
