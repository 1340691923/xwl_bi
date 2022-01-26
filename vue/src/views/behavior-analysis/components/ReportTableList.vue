<template>
  <div>
    <a-button-group>
      <a-tooltip placement="top" style="cursor: pointer">
        <template slot="title">
          <span>报表</span>
        </template>
        <a-button type="link" class="actions_xwl_btn" icon="bar-chart" @click="lookReportList" />
      </a-tooltip>
    </a-button-group>
    <a-drawer
      width="30%"
      title="已存报表"
      placement="right"
      :visible="drawerShow"
      :after-visible-change="afterVisibleChange"
      @close="onClose"
    >
      <div
        style="display: flex;width: 100%;padding: 5px;height: 100%;flex-direction: column;justify-content: flex-start;align-items: flex-start"
      >
        <el-select v-model="rtType" style="margin: 0px auto;width: 90%" placeholder="全部报表" @change="getReportTableList">
          <el-option
            v-for="(v,k,index) in tableTypeMap"
            :key="index"
            :label="v"
            :value="Number(k)"
          />
        </el-select>
        <el-input
          v-model="input"
          placeholder="请输入报表名称"
          prefix-icon="el-icon-search"
          style="margin: 0px auto;margin-top: 20px;width: 90%"
          @input="fliter"
        />
        <div
          v-infinite-scroll="handleInfiniteOnLoad"
          style="margin: 0px auto;margin-top: 20px;width: 90%"
          class="infinite-container"
          :infinite-scroll-disabled="busy"

          :infinite-scroll-distance="1000"
        >
          <a-list :data-source="data" size="small" :bordered="true" :split="true">
            <a-list-item slot="renderItem" slot-scope="item, index">

              <a-tooltip placement="left" style="cursor: pointer" @click="toTheReportTable(item.rt_type,item.id)">
                <template slot="title">
                  <span>{{ item.remark }}</span>
                </template>
                <a-list-item-meta :description="tableTypeMap[item.rt_type]">
                  <a slot="title">{{ item.name }}</a>
                </a-list-item-meta>
              </a-tooltip>

              <div>
                <el-button type="warning" icon="el-icon-delete" circle @click="deleteByID(item.id,item.name)" />
              </div>
            </a-list-item>
            <div v-if="loading && !busy" class="loading-container">
              <a-spin />
            </div>
          </a-list>
        </div>
      </div>
    </a-drawer>
  </div>
</template>
<script>

import infiniteScroll from 'vue-infinite-scroll'
import { DeleteReportTableByID, ReportTableList } from '@/api/pannel'
import { filterData } from '@/utils/table'

export default {
  directives: { infiniteScroll },
  props: {
    rtType: {
      type: Number,
      default: 0
    }
  },
  data() {
    return {
      input: '',
      tableData: [],
      rtType: this.rt_type,
      drawerShow: false,
      data: [],
      loading: false,
      busy: false,
      tableTypeMap: {
        0: '全部报表',
        1: '事件分析',
        2: '留存分析',
        3: '漏斗分析',
        4: '智能路径分析',
        5: '用户属性分析'
      },
      routerMap: {
        1: '/behavior-analysis/event/',
        2: '/behavior-analysis/retention/',
        3: '/behavior-analysis/funnel/',
        4: '/behavior-analysis/funnel/',
        5: '/user-analysis/attr/'
      }
    }
  },
  beforeMount() {
    this.getReportTableList()
  },
  methods: {
    toTheReportTable(rt, id) {
      this.$router.push({ path: this.routerMap[rt] + id })
    },
    deleteByID(id, name) {
      this.$confirm('确定删除名字为 【' + name + '】 的报表吗?', '警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await DeleteReportTableByID({ appid: this.$store.state.baseData.EsConnectID, id: id })
          if (res.code != 0) {
            this.$message({
              offset: 60,
              type: 'error',
              message: res.msg
            })
            return
          } else {
            this.$message({
              offset: 60,
              type: 'success',
              message: res.msg
            })
            this.getReportTableList()
          }
        })
        .catch(err => {
          console.error(err)
        })
    },
    fliter() {
      this.data = filterData(this.tableData, this.input)
    },
    afterVisibleChange(val) {

    },
    onClose() {
      this.drawerShow = false
    },
    lookReportList() {
      this.getReportTableList()
      this.drawerShow = true
    },
    async getReportTableList() {
      const res = await ReportTableList({ appid: this.$store.state.baseData.EsConnectID, rt_type: this.rtType })
      if (res.code != 0) {
        this.$message({
          offset: 60,
          type: 'error',
          message: res.msg
        })
        return
      }
      if (res.data == null) {
        res.data = []
      }
      this.tableData = res.data
      this.fliter()
    },
    handleInfiniteOnLoad() {

    }
  }
}
</script>
<style>
.infinite-container {
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  overflow: auto;
  padding: 8px 24px;
  height: 700px;
}

.loading-container {
  position: absolute;
  bottom: 40px;
  width: 100%;
  text-align: center;
}

.actions_xwl_btn:hover {
  color: orangered;
}
</style>
