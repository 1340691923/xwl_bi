<template>
  <div style="width: 100%;height: 100%;">

    <el-card shadow="hover" style="width: 100%">
      <div slot="header" style="display: flex; align-items: center; justify-content: space-between;">
        <span class="echartBox_title" @click="toRtPage">{{ name }}(用户属性)</span>
        <div>
          <a-tooltip placement="top" style="cursor: pointer">
            <template slot="title">
              <span>拖移报表</span>
            </template>
            <el-button type="warning" class="drageTag" icon="el-icon-rank" circle />
          </a-tooltip>
          <a-tooltip placement="top" style="cursor: pointer">
            <template slot="title">
              <span>切换成表格或者图表</span>
            </template>
            <el-button type="primary" :icon="totable?'el-icon-s-data':'el-icon-s-grid'" circle @click="changeViewType" />
          </a-tooltip>
          <a-tooltip placement="top" style="cursor: pointer">
            <template slot="title">
              <span>刷新</span>
            </template>
            <el-button icon="el-icon-refresh" circle @click="go" />
          </a-tooltip>
          <a-tooltip placement="top" style="cursor: pointer">
            <template slot="title">
              <span>下载表数据</span>
            </template>
            <el-button type="success" icon="el-icon-download" circle @click="download" />
          </a-tooltip>
        </div>
      </div>
      <a-spin tip="计算中..." :spinning="spinning">
        <div class="spin-content">

          <user-attr-result
            v-if="userAttrResShow"
            :ref="getRef"
            :tablelabel="getTableLabel"
            :class-name="getRef"
            :show-charts="tocharts"
            :show-table="totable"
            empty-text="暂无结果，请调整查询条件"
            :user-attr-res="userAttrRes"
            @go="go"
          />

        </div>
      </a-spin>
    </el-card>
  </div>
</template>

<script>

import { GetConfigs, UserAttrList } from '@/api/analysis'

export default {
  name: 'Trace',
  components: {
    'UserAttrResult': () => import('@/views/behavior-analysis/components/UserAttrResult')
  },
  props: {
    name: {
      type: String,
      default: ''
    },
    data: {
      type: Object,
      default: {}
    },
    id: {
      type: String,
      default: ''
    },
    filterDate: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      userAttrRes: [],
      spinning: false,
      allAttrOptions: [],
      form: JSON.parse(this.data),
      tableHeader: [],
      groupData: [],
      traceRes: [],
      userAttrResShow: true,
      prevCount: 0,
      tocharts: true,
      totable: false
    }
  },
  computed: {
    getRef() {
      return 'userAttr' + this.id
    },
    getTableLabel() {
      let tmp = ''
      switch (this.form.zhibiaoArr[1]) {
        case '1':
          tmp = '用户数'
          break
        case '2':
          tmp = '总和'
          break
        case '3':
          tmp = '均值'
          break
        case '4':
          tmp = '人均值'
          break
        case '5':
          tmp = '中位值'
          break
        case '6':
          tmp = '最大值'
          break
        case '7':
          tmp = '最小值'
        case '8':
          tmp = '去重值'
          break
      }
      return this.form.zhibiaoArr[0] + '.' + tmp
    }
  },

  watch: {},
  beforeMount() {
    this.init()
  },
  methods: {
    download() {
      this.$refs[this.getRef].download(this.name)
    },
    async init() {
      await this.getMetaEventList()
      await this.go()
    },
    toRtPage() {
      this.$router.push({ path: '/user-analysis/attr/' + this.id })
    },
    changeViewType() {
      this.tocharts = !this.tocharts
      this.totable = !this.totable
      this.userAttrResShow = false
      this.$nextTick(() => {
        this.userAttrResShow = true
      })
    },
    async getMetaEventList() {
      const res = await GetConfigs({ 'appid': this.$store.state.baseData.EsConnectID })

      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })
        return
      }

      const attributeMap = res.data.attributeMap

      const eventData = { label: '事件', options: [] }
      if (attributeMap.hasOwnProperty('2')) {
        for (const v of attributeMap['2']) {
          eventData.options.push({
            value: v.attribute_name,
            label: v.show_name == '' ? v.attribute_name : v.show_name
          })
        }
      }

      this.allAttrOptions = eventData
    },

    async go() {
      this.spinning = true
      this.userAttrResShow = false
      const form = this.form
      form['appid'] = this.$store.state.baseData.EsConnectID
      const res = await UserAttrList(form)
      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })
        this.userAttrRes = []
      } else {
        if (res.data.chartRes == null) {
          res.data.chartRes = []
        }
        if (res.data.tableRes == null) {
          res.data.tableRes = []
        }
        this.userAttrRes = res.data.tableRes
      }

      this.spinning = false
      this.$nextTick(() => {
        this.userAttrResShow = true
      })
    }
  }
}
</script>

<style scoped>

.spin-content {
  min-height: 500px;
}
</style>
