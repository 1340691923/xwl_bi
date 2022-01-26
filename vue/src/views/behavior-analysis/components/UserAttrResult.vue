<template>
  <div class="right_res">
    <div>
      <div
        class="app-container"
        style="height: 100%;  background: white;"
      >

        <template>

          <div v-if="showTable">
            <div v-if="tableData.length > 0">
              <div style="display: flex; align-items: center; justify-content: space-between;">

                <div class="echartBox_title" />
                <div class="echartBox_title">

                  <a-button-group>

                    <a-tooltip placement="top" style="cursor: pointer">
                      <template slot="title">
                        <span>柱状图</span>
                      </template>
                      <a-button
                        :type="chartType ==2?'primary':'default'"
                        icon="bar-chart"
                        @click.native="chartType=2"
                      />
                    </a-tooltip>
                    <a-tooltip placement="top" style="cursor: pointer">
                      <template slot="title">
                        <span>饼状图</span>
                      </template>
                      <a-button
                        :type="chartType ==3?'primary':'default'"
                        icon="pie-chart"
                        @click.native="chartType=3"
                      />
                    </a-tooltip>
                  </a-button-group>
                </div>
              </div>

              <user-attr-bar v-if="g2Show && chartType == 2" :show-label="tablelabel" :chart-data="tableData" />
              <user-attr-bar2 v-if="g2Show && chartType == 3" :show-label="tablelabel" :chart-data="tableData" />
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
                style="padding: 20px"
                :input="input"
                :limit="Number(10)"
                :table-list="tableData"
                :table-info="tableInfo"
              >
                <el-table-column slot="operate" label="指标" align="center" sortable prop="name" />
                <el-table-column slot="operate" :label="tablelabel" width="200" align="center" sortable prop="value" />
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
    'UserAttrBar': () => import('@/components/Charts/UserAttrBar'),
    'UserAttrBar2': () => import('@/components/Charts/UserAttrBar2')
  },
  props: {
    tablelabel: {
      type: String,
      default: ''
    },
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
    userAttrRes: {
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
      chartType: 2,
      showList: [],
      tableData: [],
      chartData: []
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
    rewriteNodeName(event, num) {
      return `${event}_${num}`
    },
    download(fName) {
      elTable2Excel(this, 'pagetable', `用户属性分析:${fName}`)
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
      this.tableData = this.userAttrRes
    },
    init() {
      this.initTableData()
      this.refreshData()
    }
  }

}
</script>

<style scoped src="@/styles/UserAttrResult.css"/>

