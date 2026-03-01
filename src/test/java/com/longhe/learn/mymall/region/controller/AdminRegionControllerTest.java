package com.longhe.learn.mymall.region.controller;

import com.longhe.learn.mymall.core.mapper.RedisUtil;
import com.longhe.learn.mymall.core.model.ReturnNo;
import com.longhe.learn.mymall.core.util.JwtHelper;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.mock.mockito.MockBean;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;
import org.springframework.transaction.annotation.Transactional;

import static org.hamcrest.CoreMatchers.is;
import static org.junit.jupiter.api.Assertions.*;

@SpringBootTest
@AutoConfigureMockMvc
@Transactional
class AdminRegionControllerTest {
    @Autowired
    private MockMvc mockMvc;
    @MockBean
    private RedisUtil redisUtil;
    private static String adminToken;
    private static String wrongToken;
    private final String ADMIN_SUB_REGIONS = "/shops/{did}/regions/{id}/subregions";

    @BeforeAll
    static void setUp() {
        JwtHelper jwtHelper = new JwtHelper();
        adminToken = jwtHelper.createToken(1L, "13088admin", 0L, 1, 3600);
        wrongToken = jwtHelper.createToken(1L, "13088admin", 1L, 1, 3600);
    }

    @Test
    void createSubRegions() throws Exception {
        Mockito.when(redisUtil.hasKey(Mockito.anyString())).thenReturn(false);
        Mockito.when(redisUtil.get(Mockito.anyString())).thenReturn(null);
        Mockito.when(redisUtil.set(Mockito.anyString(), Mockito.any(), Mockito.anyLong())).thenReturn(true);

        String body = "{\"name\":\"DongCheng\", \"shortName\":\"DongCheng\", \"mergerName\":\"Beijing DongCheng\",\"pinyin\":\"Dong Cheng\",\"lng\":\"116.416357\", \"lat\":\"39.928353\",\"areaCode\":\"110101000000\",\"zipCode\":\"00100000\",\"cityCode\":\"010\"}";

        this.mockMvc.perform(MockMvcRequestBuilders.post(ADMIN_SUB_REGIONS, 0, 16)
                        .header("authorization", adminToken)
                        .contentType(MediaType.APPLICATION_JSON_VALUE)
                        .content(body))
                .andExpect(MockMvcResultMatchers.status().isCreated())
                .andExpect(MockMvcResultMatchers.jsonPath("$.errno", is(ReturnNo.CREATED.getErrNo())))
                        .andExpect(MockMvcResultMatchers.jsonPath("$.data.name", is("DongCheng")))
                .andExpect(MockMvcResultMatchers.jsonPath("$.data.id").exists());
    }

    @Test
    void createSubRegionsWrongToken() throws Exception {
        Mockito.when(redisUtil.hasKey(Mockito.anyString())).thenReturn(false);
        Mockito.when(redisUtil.get(Mockito.anyString())).thenReturn(null);
        Mockito.when(redisUtil.set(Mockito.anyString(), Mockito.any(), Mockito.anyLong())).thenReturn(true);

        String body = "{\"name\":\"DongCheng\", \"shortName\":\"DongCheng\", \"mergerName\":\"Beijing DongCheng\",\"pinyin\":\"Dong Cheng\",\"lng\":\"116.416357\", \"lat\":\"39.928353\",\"areaCode\":\"110101000000\",\"zipCode\":\"00100000\",\"cityCode\":\"010\"}";

        this.mockMvc.perform(MockMvcRequestBuilders.post(ADMIN_SUB_REGIONS, 0, 16)
                        .header("authorization", wrongToken)
                        .contentType(MediaType.APPLICATION_JSON_VALUE)
                        .content(body))
                .andExpect(MockMvcResultMatchers.status().isForbidden())
                .andExpect(MockMvcResultMatchers.jsonPath("$.errno", is(ReturnNo.RESOURCE_ID_OUTSCOPE.getErrNo())));
    }

    @Test
    void createSubRegionsGivenNotPlatform() throws Exception {
        Mockito.when(redisUtil.get(Mockito.anyString())).thenReturn(null);
        Mockito.when(redisUtil.set(Mockito.anyString(), Mockito.any(), Mockito.anyLong())).thenReturn(true);

        String body = "{\"name\":\"东城区风云再起\", \"shortName\":\"风云再起\", \"mergerName\":\"北京，东城，风云再起\",\"pinyin\":\"FengYunZaiQi\",\"lng\":\"116.416357\", \"lat\":\"39.928353\",\"areaCode\":\"110101000000\",\"zipCode\":\"00100000\",\"cityCode\":\"010\"}";

        this.mockMvc.perform(MockMvcRequestBuilders.post(ADMIN_SUB_REGIONS, 1, 16)
                        .header("authorization", adminToken)
                        .contentType(MediaType.APPLICATION_JSON_VALUE)
                        .content(body))
                .andExpect(MockMvcResultMatchers.status().isForbidden())
                .andExpect(MockMvcResultMatchers.jsonPath("$.errno", is(ReturnNo.RESOURCE_ID_OUTSCOPE.getErrNo())));
    }

    @Test
    void createSubRegionGivenWrongFormat() throws Exception {
        Mockito.when(redisUtil.get(Mockito.anyString())).thenReturn(null);
        Mockito.when(redisUtil.set(Mockito.anyString(), Mockito.any(), Mockito.anyLong())).thenReturn(true);

        String body = "{\"name\":\"东城区风云再起\", \"shortName\":\"风云再起\", \"mergerName\":\"北京，东城，风云再起\",\"pinyin\":\"FengYunZaiQi\",\"lng\":\"116.416357\", \"lat\":\"s39.928353\",\"areaCode\":\"110101000000\",\"zipCode\":\"00100000\",\"cityCode\":\"010\"}";

        this.mockMvc.perform(MockMvcRequestBuilders.post(ADMIN_SUB_REGIONS, 0, 16)
                        .header("authorization", adminToken)
                        .contentType(MediaType.APPLICATION_JSON_VALUE)
                        .content(body))
                .andExpect(MockMvcResultMatchers.status().isBadRequest())
                .andExpect(MockMvcResultMatchers.jsonPath("$.errno", is(ReturnNo.FIELD_NOTVALID.getErrNo())));
    }

    @Test
    void createSubRegionGivenNullString() throws Exception {
        Mockito.when(redisUtil.get(Mockito.anyString())).thenReturn(null);
        Mockito.when(redisUtil.set(Mockito.anyString(), Mockito.any(), Mockito.anyLong())).thenReturn(true);

        String body = "{\"name\":\"东城区风云再起\", \"shortName\":\"风云再起\", \"mergerName\":\"北京，东城，风云再起\",\"pinyin\":\"FengYunZaiQi\",\"lng\":\"\", \"lat\":\"39.928353\",\"areaCode\":\"110101000000\",\"zipCode\":\"00100000\",\"cityCode\":\"010\"}";

        this.mockMvc.perform(MockMvcRequestBuilders.post(ADMIN_SUB_REGIONS, 0, 16)
                        .header("authorization", adminToken)
                        .contentType(MediaType.APPLICATION_JSON_VALUE)
                        .content(body))
                .andExpect(MockMvcResultMatchers.status().isBadRequest())
                .andExpect(MockMvcResultMatchers.jsonPath("$.errno", is(ReturnNo.FIELD_NOTVALID.getErrNo())));
    }
}