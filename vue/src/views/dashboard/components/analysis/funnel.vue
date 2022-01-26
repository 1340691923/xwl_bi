<template>
  <div style="width: 100%;height: 100%;">

    <el-card shadow="hover" style="width: 100%">
      <div slot="header" style="display: flex; align-items: center; justify-content: space-between;">
        <span class="echartBox_title" @click="toRtPage">{{ name }}(漏斗)</span>
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
          <funnel-result
            v-if="funnelResShow"
            :ref="getRef"
            v-model="form.date"
            empty-text="暂无结果，请调整查询条件"
            :class-name="getRef"
            :show-charts="tocharts"
            :show-table="totable"
            :table-header="tableHeader"
            :group-data="groupData"
            :funnel-res="funnelRes"
            @go="go"
          />
        </div>
      </a-spin>
    </el-card>
  </div>
</template>

<script>
import { FunnelList, GetConfigs } from '@/api/analysis'
import moment from 'moment'

export default {
  name: 'Funnel',
  components: {
    'FunnelResult': () => import('@/views/behavior-analysis/components/FunnelResult')
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
      spinning: false,
      allAttrOptions: [],
      form: JSON.parse(this.data),
      tableHeader: [],
      groupData: [],
      funnelRes: [],
      funnelResShow: true,
      prevCount: 0,
      tocharts: true,
      totable: false
    }
  },
  computed: {
    getRef() {
      return 'funnelRes' + this.id
    }
  },

  watch: {
    'filterDate': {
      immediate: true,
      handler() {
        if (this.filterDate.length > 0) {
          this.form.date = this.filterDate
          Vue.set(this.form, 'date', this.filterDate)
          this.go()
        }
      }
    }
  },
  beforeMount() {
    if (this.filterDate.length == 0) {
      this.form.date = [
        moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD'),
        moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD')
      ]
    }
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
      this.$router.push({ path: '/behavior-analysis/funnel/' + this.id })
    },
    changeViewType() {
      this.tocharts = !this.tocharts
      this.totable = !this.totable
      this.funnelResShow = false
      this.$nextTick(() => {
        this.funnelResShow = true
      })
    },

    getWashData(data) {
      const funnelRes = []
      let maxCount = 0

      if (data != null) {
        maxCount = Math.max.apply(Math, data.map(function(o) {
          return o.count
        }))
      }
      this.prevCount = 0

      for (const i in this.form.zhibiaoArr) {
        const index = Number(i) + 1
        var obj = {
          level_desc: this.form.zhibiaoArr[i].eventNameDisplay,
          level_index: index,
          showTitle: `步骤${index}:${this.form.zhibiaoArr[i].eventNameDisplay}`,
          count: 0,
          conversionScale: 0,
          washScale: 100,
          lostScale: 100,
          succScale: 0,
          max: 100
        }

        if (data != null) {
          for (const v of data) {
            if (v.level_index == index) {
              if (index == 1) {
                obj.count = v.count
                obj.succScale = 100.00
                obj.lostScale = 0
                obj.ui = v.ui
              } else {
                obj.count = v.count
                obj.ui = v.ui
                obj.succScale = (Number((v.count / this.prevCount)) * 100).toFixed(2)
                obj.lostScale = (100 - obj.succScale).toFixed(2)
              }

              this.prevCount = v.count
              continue
            }
          }
        }
        funnelRes.push(obj)
      }

      for (const i in funnelRes) {
        funnelRes[i].conversionScale = (Number((funnelRes[i].count / funnelRes[0].count)) * 100).toFixed(2)
        funnelRes[i].washScale = (100 - parseFloat(funnelRes[i].conversionScale)).toFixed(2)

        if (i == Number(funnelRes.length) - 1) {
          continue
        }

        funnelRes[i].succScale = funnelRes[Number(i) + 1].succScale
        funnelRes[i].lostScale = funnelRes[Number(i) + 1].lostScale
      }
      return funnelRes
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

    getGroupByLable(v) {
      for (const option of this.allAttrOptions.options) {
        if (option.value == v) {
          return option.label
        }
      }
    },

    async go() {
      this.spinning = true
      this.funnelResShow = false
      const form = this.form
      form['appid'] = this.$store.state.baseData.EsConnectID
      const res = await FunnelList(form)
      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })
        this.funnelRes = []
        this.groupData = []

        this.spinning = false
        this.$nextTick(() => {
          this.funnelResShow = true
        })
        return
      }

      this.funnelRes = this.getWashData(res.data.groupData['总体'])

      for (const k in res.data.groupData) {
        res.data.groupData[k] = this.getWashData(res.data.groupData[k])
      }

      this.groupData = res.data.groupData

      this.tableHeader = [
        this.form.groupBy.length > 0 ? this.getGroupByLable(this.form.groupBy[0]) : '总体'
      ]
      for (const v of this.funnelRes) {
        this.tableHeader.push(v.showTitle)
      }

      this.spinning = false
      this.$nextTick(() => {
        this.funnelResShow = true
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
