package com.longhe.learn.mymall.region.mapper;

import com.longhe.learn.mymall.region.mapper.po.RegionPo;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.transaction.annotation.Transactional;

import java.util.Optional;

import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest
@Transactional
class RegionPoMapperTest {

    @Autowired
    private RegionPoMapper regionPoMapper;

}