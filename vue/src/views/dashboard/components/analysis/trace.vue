<template>
  <div style="width: 100%;height: 100%;">

    <el-card shadow="hover" style="width: 100%">
      <div slot="header" style="display: flex; align-items: center; justify-content: space-between;">
        <span class="echartBox_title" @click="toRtPage">{{ name }}(智能路径)</span>
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

          <trace-result
            v-if="traceResShow"
            :ref="getRef"
            v-model="form.date"
            :class-name="getRef"
            :show-charts="tocharts"
            :show-table="totable"
            empty-text="暂无结果，请调整查询条件"
            :trace-charts-res="traceChartsRes"
            :trace-table-res="traceTableRes"
            @go="go"
          />

        </div>
      </a-spin>
    </el-card>
  </div>
</template>

<script>
import { GetConfigs, TraceList } from '@/api/analysis'
import moment from 'moment'

export default {
  name: 'Trace',
  components: {
    'TraceResult': () => import('@/views/behavior-analysis/components/TraceResult')
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
      traceChartsRes: [],
      traceTableRes: [],
      spinning: false,
      allAttrOptions: [],
      form: JSON.parse(this.data),
      tableHeader: [],
      groupData: [],
      traceRes: [],
      traceResShow: true,
      prevCount: 0,
      tocharts: true,
      totable: false
    }
  },
  computed: {
    getRef() {
      return 'traceRes' + this.id
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
      this.$router.push({ path: '/behavior-analysis/trace/' + this.id })
    },
    changeViewType() {
      this.tocharts = !this.tocharts
      this.totable = !this.totable
      this.traceResShow = false
      this.$nextTick(() => {
        this.traceResShow = true
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
      this.traceResShow = false
      const form = this.form
      form['appid'] = this.$store.state.baseData.EsConnectID
      const res = await TraceList(form)
      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })
        this.traceChartsRes = []
        this.traceTableRes = []
      } else {
        if (res.data.chartRes == null) {
          res.data.chartRes = []
        }
        if (res.data.tableRes == null) {
          res.data.tableRes = []
        }
        this.traceChartsRes = res.data.chartRes
        this.traceTableRes = res.data.tableRes
      }

      this.spinning = false
      this.$nextTick(() => {
        this.traceResShow = true
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
