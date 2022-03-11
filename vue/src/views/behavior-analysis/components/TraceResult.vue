<template>
  <div class="right_res">
    <div>
      <div
        class="app-container"
        style="height: 100%;  background: white;"
      >
        <div style="display: flex; align-items: center; justify-content: space-between;">
          <div class="echartBox_title">
            <date v-if="dateShow" v-model="filterDate" @changeDate="filterDateCall" />
          </div>

        </div>
        <template>

          <div v-if="showTable">

            <div v-if="tableData.length > 0">
              <trace-sankey v-if="g2Show" :chart-data="traceChartsRes" />
            </div>

          </div>
          <div v-if="showCharts">
            <div
              v-if="tableData.length <= 0"
              style="background: white !important;padding: 40px;width: 300px;height: 300px; text-align: center;margin: 0px auto"
            >
              <a-empty>
                <span slot="description">{{ emptyText }}</span>
              </a-empty>
            </div>
            <div v-else style="margin-top: 20px">
              <el-input v-model="input" class="filter-item" placeholder="输入关键字进行过滤" clearable style="width: 300px" />
              <page-table
                v-if="tableShow"
                ref="pagetable"
                style="margin-top: 20px"
                :input="input"
                :limit="Number(10)"
                :table-list="tableData"
                :table-info="tableInfo"
              >
                <el-table-column slot="operate" label="路径" align="center" sortable prop="trace" />

                <el-table-column slot="operate" label="人数" width="100" align="center" sortable prop="user_count">
                  <template slot-scope="scope">
                    <a style="color: #6bb8ff" @click="drillDown(scope.row.ui)">&nbsp;&nbsp;&nbsp;&nbsp;{{ scope.row.user_count }}</a>
                    <add-user-group :uid="scope.row.ui" />
                  </template>
                </el-table-column>

                <el-table-column slot="operate" label="比例" width="100" align="center" sortable>
                  <template slot-scope="scope">
                    {{ (Number((scope.row.user_count / allUserNum)) * 100).toFixed(2) }}%
                  </template>
                </el-table-column>

              </page-table>
            </div>
          </div>

        </template>
      </div>
    </div>

  </div>
</template>

<script>
import { elTable2Excel } from '@/utils/download'

export default {
  name: 'RetentionResult',
  components: {
    'PageTable': () => import('@/components/PageTable'),
    'Date': () => import('@/components/AnalyseTools/FilterDate/Date'),
    'TraceSankey': () => import('@/components/Charts/TraceSankey'),
    'AddUserGroup': () => import('@/views/behavior-analysis/components/AddUserGroup')
  },
  props: {
    className: {
      type: String,
      default: 'charts'
    },
    emptyText: {
      type: String,
      default: '选择完分析条件后，请点击“计算”'
    },
    value: {
      type: Array,
      default: []
    },
    traceTableRes: {
      type: Array,
      default: []
    },
    traceChartsRes: {
      type: Array,
      default: []
    },
    showCharts: {
      type: Boolean,
      default: true
    },
    showTable: {
      type: Boolean,
      default: true
    }
  },
  data() {
    return {
      dateShow: true,
      input: '',
      tableHeaderShow: [],
      groupDataShow: {},
      g2Show: true,
      tableShow: true,
      filterDate: this.value,
      tableInfo: [{ slot: 'operate' }],
      chartType: 1,
      showList: [],
      tableData: [],
      chartData: [],
      allUserNum: 0
    }
  },
  computed: {},
  watch: {

    value: {
      deep: true,
      handler() {
        this.dateShow = false
        this.$nextTick(() => {
          this.dateShow = true
        })
      }
    }
  },
  beforeMount() {
    this.init()
  },
  methods: {
    drillDown(ui) {
      console.log('ui', ui)
      this.$store.dispatch('baseData/SETUI', ui)
      this.$router.push({ path: '/user-analysis/user_list' })
    },
    rewriteNodeName(event, num) {
      return `${event}_${num}`
    },
    download(fName) {
      elTable2Excel(this, 'pagetable', `智能路径分析:${fName}`)
    },
    refreshData() {
      this.g2Show = false
      this.tableShow = false
      this.$nextTick(() => {
        this.g2Show = true
        this.tableShow = true
      })
    },
    filterDateCall(date) {
      this.filterDate = date
      this.$emit('input', this.filterDate)
      this.$emit('go')
    },
    initTableData() {
      this.allUserNum = 0
      for (const v of this.traceTableRes) {
        this.allUserNum = this.allUserNum + v.user_count
      }
      this.tableData = this.traceTableRes
    },
    initChartData() {
      this.chartData = []
      if (this.traceChartsRes.length <= 0) {
        return
      }
      const eventArr = []
      const targetArr = []
      const eventSet = new Map()

      for (const k in this.traceChartsRes) {
        const traceCharts = this.traceChartsRes[k]

        eventSet.set(traceCharts['event'][0], 1)
        eventSet.set(traceCharts['event'][1], 1)

        targetArr.push({
          source: traceCharts['event'][0],
          target: traceCharts['event'][1],
          value: traceCharts['sum_user_count']
        })
      }
      eventSet.forEach((v, k, tmp) => {
        const obj = {
          name: k
        }
        eventArr.push(obj)
      })
      this.chartData.push(eventArr)
      this.chartData.push(targetArr)
    },
    init() {
      this.initTableData()
      this.initChartData()
      this.refreshData()
    }
  }

}
</script>

<style scoped src="@/styles/trace-res.css"/>

