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
            <a-divider type="vertical" />
            <a-cascader
              v-model="step"
              style="width: 160px"
              :options="stepOptions"
              placeholder="请选择步骤"
              @change="changeStep"
            />

          </div>
          <div class="echartBox_title">
            <!--<a-button-group>
              <a-tooltip placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>转化图</span>
                </template>
                <a-button
                  :type="chartType ==1?'primary':'default'"
                  icon="bar-chart"
                  @click.native="chartType=1"
                />
              </a-tooltip>
            </a-button-group>-->
          </div>
        </div>
        <template v-if="showCharts">
          <funnel-bar v-if="g2Show && showList.length > 0 " :class-name="className" :chart-data="showList" />
          <div
            v-else
            style="background: white !important;padding: 40px;width: 300px;height: 300px; text-align: center;margin: 0px auto"
          >
            <a-empty>
              <span slot="description">{{ emptyText }}</span>
            </a-empty>
          </div>
        </template>
      </div>
      <template v-if="showTable">
        <div v-if="g2Show" style="background: white">
          <div class="filter-container" style="background: white;padding: 20px">
            <el-radio-group v-model="tableType" class="filter-item">
              <el-radio-button label="conversion">转化</el-radio-button>
              <el-radio-button label="drain">流失</el-radio-button>
            </el-radio-group>
            <el-input v-model="input" class="filter-item" placeholder="输入关键字进行过滤" style="width: 300px" />
          </div>
          <div v-if="tableTitle.toString() !=''" style="width: 100%;padding: 20px;background: white">
            <div
              style="width: 300px;text-align: center;margin: 0px auto;font-weight: bolder;font-size: 20px;color: #909399"
            >
              {{ tableTitle.toString() }}
            </div>
          </div>

          <page-table
            v-if="tableShow"
            ref="pagetable"
            table-ref="table"
            style="padding: 20px"
            :input="input"
            :show-title="tableTitle.toString()"
            :table-list="tableData"
            :table-info="tableInfo"
          >
            <el-table-column slot="operate" :label="tableHeaderShow[0]" align="center" sortable prop="groupKey" />
            <el-table-column
              v-for="(v,k,index) in tableHeaderShow.slice(1,tableHeaderShow.length)"
              slot="operate"
              :key="index"
              :label="v"
              align="center"
              sortable
            >
              <template slot-scope="scope">
                <div
                  style="display: flex;flex-direction: column;align-items: center;justify-content: center;width: 73%;max-width: 92px;margin: auto;"
                >
                  <div>
                    <a style="color: #6bb8ff" @click="drillDown(scope.row.countArr[k].ui)">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;{{
                      scope.row.countArr[k].count
                    }}&nbsp;&nbsp;</a>
                    <add-user-group :uid="scope.row.countArr[k].ui" />
                  </div>
                  <div style="text-align: center;">
                    <template v-if="tableType == 'conversion'">{{
                      scope.row.countArr[k].conversionScale
                    }}%&nbsp;&nbsp;
                    </template>
                    <template v-else>{{ scope.row.countArr[k].washScale }}%&nbsp;&nbsp;</template>
                  </div>
                </div>
              </template>
            </el-table-column>
          </page-table>
        </div>
      </template>
    </div>

  </div>
</template>

<script>

import { elTable2Excel } from '@/utils/download'

export default {
  name: 'FunnelResult',
  components: {
    'FunnelBar': () => import('@/components/Charts/FunnelBar'),
    'PageTable': () => import('@/components/PageTable'),
    'Date': () => import('@/components/AnalyseTools/FilterDate/Date'),
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
    funnelRes: {
      type: Array,
      default: []
    },
    groupData: {
      type: Object,
      default: {}
    },
    tableHeader: {
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
      step: [-1],
      stepOptions: [],
      tableType: 'conversion',
      tableInfo: [{ slot: 'operate' }],
      chartType: 0,
      showList: [],
      tableData: []
    }
  },
  computed: {
    tableTitle() {
      return `全步骤（共${this.showList.length}步）的用户${this.tableType == 'conversion' ? '转化率' : '流失率'}`
    }
  },
  watch: {
    funnelRes: {
      deep: true,
      handler() {
        this.init()
        this.changeStep()
      }
    },
    groupData: {
      deep: true,
      handler() {
        this.init()
        this.changeStep()
      }
    },
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
  mounted() {
    this.init()
    this.changeStep()
  },
  methods: {
    drillDown(ui) {
      this.$store.dispatch('baseData/SETUI', ui)
      this.$router.push({ path: '/user-analysis/user_list' })
    },
    download(fName) {
      elTable2Excel(this, 'pagetable', `漏斗分析:${fName}`)
    },
    NaN2Zero(v) {
      if (isNaN(v)) {
        return 0
      }
      return v
    },
    changeTable() {
      const tableData = []

      for (const k in this.groupDataShow) {
        const countArr = []
        for (const v of this.groupDataShow[k]) {
          countArr.push({
            count: v.count,
            ui: v.ui,
            conversionScale: this.NaN2Zero(v.conversionScale),
            washScale: this.NaN2Zero(v.washScale)
          })
        }
        tableData.push({ 'groupKey': k, countArr: countArr })
      }

      this.tableData = tableData
    },
    changeStep() {
      switch (this.step.length) {
        case 0:
          this.step = [-1]
        case 1:
          this.showList = JSON.parse(JSON.stringify(this.funnelRes))
          this.groupDataShow = JSON.parse(JSON.stringify(this.groupData))
          this.tableHeaderShow = JSON.parse(JSON.stringify(this.tableHeader))

          break
        case 2:
          const stepStart = this.step[0]
          const stepOver = this.step[1]

          const funnelRes = []

          for (const v of this.funnelRes) {
            if (v.level_index >= stepStart && v.level_index <= stepOver) {
              funnelRes.push(v)
            }
          }
          const groupDataShow = {}
          for (const k in this.groupData) {
            const tmp = []

            for (const v of this.groupData[k]) {
              if (v.level_index >= stepStart && v.level_index <= stepOver) {
                tmp.push(v)
              }
            }
            groupDataShow[k] = tmp
          }
          this.groupDataShow = groupDataShow
          this.showList = funnelRes
          const groupTitle = this.tableHeader[0]
          const tableHeader = [groupTitle]
          for (const k in this.tableHeader) {
            if (k >= stepStart && k <= stepOver) {
              tableHeader.push(this.tableHeader[k])
            }
          }

          this.tableHeaderShow = tableHeader
      }
      this.changeTable()
      this.refreshData()
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
    init() {
      const stepOpt = [
        {
          value: -1,
          label: '全步骤'
        }
      ]
      console.log(this.funnelRes)
      for (const i in this.funnelRes) {
        if (i == this.funnelRes.length - 1) {
          break
        }

        const stepNum = Number(i) + 1
        var obj = {
          value: stepNum,
          label: '步骤' + stepNum,
          children: []
        }
        let childrenStep = Number(stepNum) + 1

        for (const k in [...new Array(this.funnelRes.length - stepNum).keys()]) {
          obj.children.push({
            value: childrenStep,
            label: '步骤' + childrenStep
          })
          childrenStep++
        }
        stepOpt.push(obj)
      }
      this.stepOptions = stepOpt
    }
  }
}
</script>

<style scoped src="@/styles/funnel-res.css"/>
