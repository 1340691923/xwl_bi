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
            <el-select v-model="windowTimeFormat2" size="mini" style="width: 100px">
              <el-option
                v-for="item in windowTimeOpt"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>

          </div>
          <div class="echartBox_title">

            <a-button-group v-if="chartType == 2 || chartType == 3">
              <a-tooltip placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>比例</span>
                </template>
                <a-button
                  :type="isScale == true?'primary':'default'"
                  @click.native="isScale=true"
                >比例
                </a-button>
              </a-tooltip>
              <a-tooltip placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>人数</span>
                </template>
                <a-button
                  :type="isScale == false?'primary':'default'"

                  @click.native="isScale=false"
                >人数
                </a-button>
              </a-tooltip>
            </a-button-group>

            <!--<a-button-group>
              <a-tooltip placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>趋势图</span>
                </template>
                <a-button
                  :type="chartType ==1?'primary':'default'"
                  icon="table"
                  @click.native="chartType=1"
                />
              </a-tooltip>
              <a-tooltip placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>堆积图</span>
                </template>
                <a-button
                  :type="chartType ==2?'primary':'default'"
                  icon="fund"
                  @click.native="chartType=2"
                />
              </a-tooltip>
              <a-tooltip placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>累计图</span>
                </template>
                <a-button
                  :type="chartType ==3?'primary':'default'"
                  icon="fund"
                  @click.native="chartType=3"
                />
              </a-tooltip>
              <a-tooltip placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>分布图</span>
                </template>
                <a-button
                  :type="chartType ==4?'primary':'default'"
                  icon="fund"
                  @click.native="chartType=4"
                />
              </a-tooltip>
            </a-button-group>-->
          </div>
        </div>
        <template>
          <div
            v-if="tableData.length <= 0"
            style="background: white !important;padding: 40px;width: 300px;height: 300px; text-align: center;margin: 0px auto"
          >
            <a-empty>
              <span slot="description">{{ emptyText }}</span>
            </a-empty>
          </div>
          <div v-else-if="chartType==1">
            <div class="filter-container" style="background: white;padding: 20px">
              <el-input v-model="input" class="filter-item" placeholder="输入关键字进行过滤" style="width: 300px" />
              <el-select v-model="groupType" class="filter-item" @change="init">
                <el-option label="按日期分组" :value="Number(1)" />
                <el-option label="按事件分组" :value="Number(2)" />
              </el-select>
            </div>

            <page-table
              v-if="tableShow"
              ref="pagetable"
              :span-method="spanMethod"
              style="padding: 20px"
              :input="input"
              :limit="Number(10)"
              :table-list="tableData"
              :table-info="tableInfo"
            >
              <el-table-column
                v-for="(v,k,index) in tableTitle"
                :key="index"
                slot="operate"
                :formatter="typeFormatter"
                :label="v.label"
                align="center"
                :prop="v.prop"
              />
            </page-table>
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
    'RetentionLine': () => import('@/components/Charts/RetentionLine'),
    'RetentionLine2': () => import('@/components/Charts/RetentionLine2'),
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
    eventRes: {
      type: Object,
      default: {}
    },
    windowTimeFormat: {
      type: String,
      default: '按天'
    }
  },
  data() {
    return {
      lookTypScaleMap: {
        'retention': 'conversionScaleArr',
        'lose': 'washScaleArr'
      },
      windowTimeOpt: [
        {
          value: '按天',
          label: '按天'
        },
        {
          value: '按分钟',
          label: '按分钟'
        },
        {
          value: '按小时',
          label: '按小时'
        }, {
          value: '按周',
          label: '按周'
        },
        {
          value: '按月',
          label: '按月'
        }, {
          value: '合计',
          label: '合计'
        }
      ],
      groupType: 2,
      windowTimeFormat2: this.windowTimeFormat,
      lookTyp: 'retention',
      dateShow: true,
      input: '',
      tableHeaderShow: [],
      groupDataShow: {},
      g2Show: true,
      tableShow: true,
      filterDate: this.value,
      tableInfo: [{ slot: 'operate' }],
      chartType: 1,
      isScale: true,
      showList: [],
      spanMethodFields: {},
      tableTitle: [],
      tableData: []
    }
  },
  computed: {},
  watch: {

    'windowTimeFormat2': {
      deep: true,
      handler() {
        this.$emit('changeWindowTime', this.windowTimeFormat2)
      }
    },
    windowTimeFormat(val, oldVal) {
      this.windowTimeFormat2 = val
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
  },
  methods: {
    typeFormatter(row, column, cellValue, index) {
      if (cellValue == undefined) {
        return 0
      }
      return cellValue
    },
    refreshData() {
      this.tableShow = false
      this.$nextTick(() => {
        this.tableShow = true
      })
    },
    spanMethod({ row, column, rowIndex, columnIndex }) {
      const fields = this.spanMethodFields
      const cellValue = row[column.property]

      if (cellValue != undefined && fields.hasOwnProperty(column.property)) {
        const prevRow = this.tableData[rowIndex - 1]
        let nextRow = this.tableData[rowIndex + 1]
        if (prevRow && prevRow[column.property] === cellValue) {
          return { rowspan: 0, colspan: 0 }
        } else {
          let countRowspan = 1
          while (nextRow && (nextRow[column.property] === cellValue)) {
            nextRow = this.tableData[++countRowspan + rowIndex]
          }
          if (countRowspan > 1) {
            return { rowspan: countRowspan, colspan: 1 }
          }
        }
      }
    },
    filterDateCall(date) {
      this.filterDate = date
      this.$emit('input', this.filterDate)
      this.$emit('go')
    },
    download(fName) {
      elTable2Excel(this, 'pagetable', `事件分析:${fName}`)
    },

    compareDate(property) {
      return function(a, b) {
        var value1 = a[property]
        var value2 = b[property]
        return value1.localeCompare(value2, 'zh-CN')
      }
    },
    init() {
      if (!this.eventRes.hasOwnProperty('alldata')) return

      const useGroup = this.eventRes.use_group
      const len = this.eventRes.len
      const groupby = this.eventRes.groupby
      const alldata = this.eventRes.alldata
      const eventNameDisplayArr = this.eventRes.eventNameDisplayArr
      this.tableTitle = []
      this.spanMethodFields = {}
      this.tableData = []
      const m = {}

      if (this.groupType == 1) {
        if (len == 1 && useGroup == true) {
          this.spanMethodFields['eventNameDisplay'] = 1
        } else if (len > 1) {
          for (const v of groupby) {
            this.spanMethodFields[v] = 1
          }
        }

        for (const v of groupby) {
          const tableTitleObj = {
            label: v,
            prop: v
          }
          this.tableTitle.push(tableTitleObj)
        }

        this.tableTitle.push(
          {
            label: '指标',
            prop: 'eventNameDisplay'
          }
        )

        for (const k in alldata) {
          const v = alldata[k]
          if (v.data_group.length == 0) {
            continue
          }
          const tableDataobj = {
            eventNameDisplay: v.eventNameDisplay
          }

          for (const groupV of groupby) {
            tableDataobj[groupV] = v[groupV]
          }

          for (const k2 in v.data_group) {
            const v2 = v.data_group[k2]
            const date = v2[0]
            const amount = v2[1]
            if (!m.hasOwnProperty(date)) {
              const tableTitleObj = {
                label: date,
                prop: date
              }
              this.tableTitle.push(tableTitleObj)
              m[date] = 1
            }
            tableDataobj[date] = amount
          }
          this.tableData.push(tableDataobj)
        }
      } else {
        const SPLIT = '$$$xwl$$$'
        this.spanMethodFields['date'] = 1

        this.tableTitle.push(
          {
            label: '日期',
            prop: 'date'
          }
        )

        for (const v of groupby) {
          const tableTitleObj = {
            label: v,
            prop: v
          }
          this.tableTitle.push(tableTitleObj)
        }
        for (const v of eventNameDisplayArr) {
          const tableTitleObj = {
            label: v,
            prop: v
          }
          this.tableTitle.push(tableTitleObj)
        }
        const obj = {}
        for (const k in alldata) {
          const v = alldata[k]
          if (v.data_group.length == 0) {
            continue
          }
          const eventNameDisplay = v['eventNameDisplay']
          for (const k2 in v.data_group) {
            const date = v.data_group[k2][0]
            const amount = v.data_group[k2][1]

            const groupArr = []
            groupArr.push(date)
            for (const groupbyV of groupby) {
              groupArr.push(v[groupbyV])
            }

            const groupKey = groupArr.join(SPLIT)
            let tmp = {}
            if (obj.hasOwnProperty(groupKey)) {
              tmp = obj[groupKey]
            }
            tmp[eventNameDisplay] = amount
            tmp['date'] = date
            for (const groupbyV of groupby) {
              tmp[groupbyV] = v[groupbyV]
            }

            obj[groupKey] = tmp
          }
        }

        for (const groupKey in obj) {
          const tableObj = obj[groupKey]
          this.tableData.push(tableObj)
        }
        this.tableData.sort(this.compareDate('date'))
      }

      this.refreshData()
    }
  }

}
</script>

<style scoped src="@/styles/EventResult.css"/>
