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
            <el-input-number
              v-model="form.windowTime"
              size="mini"
              controls-position="right"
              style="width: 90px"
              :min="1"
            />
            <el-select v-model="form.windowTimeFormat" size="mini" style="width: 70px">
              <el-option
                v-for="item in windowTimeOpt"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
            <a-divider type="vertical" />
            <el-select v-model="lookTyp" size="mini" style="width: 70px">
              <el-option
                v-for="item in lookTypOpt"
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

            <a-button-group>
              <a-tooltip placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>数据表</span>
                </template>
                <a-button
                  :type="chartType ==1?'primary':'default'"
                  icon="table"
                  @click.native="chartType=1"
                />
              </a-tooltip>
              <a-tooltip placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>第N日留存</span>
                </template>
                <a-button
                  :type="chartType ==2?'primary':'default'"
                  icon="fund"
                  @click.native="chartType=2"
                />
              </a-tooltip>
              <a-tooltip placement="top" style="cursor: pointer">
                <template slot="title">
                  <span>每日留存</span>
                </template>
                <a-button
                  :type="chartType ==3?'primary':'default'"
                  icon="fund"
                  @click.native="chartType=3"
                />
              </a-tooltip>
            </a-button-group>
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

            <page-table
              v-if="tableShow"
              ref="pagetable"
              style="padding: 20px"
              :input="input"
              :limit="Number(10)"
              :table-list="tableData"
              :table-info="tableInfo"
            >

              <el-table-column
                slot="operate"
                :label="tableHeaderShow[0]"
                align="center"
                width="100"
                sortable
                prop="dates"
              />

              <el-table-column slot="operate" :key="index" :label="tableHeaderShow[1]" align="center" sortable>
                <template slot-scope="scope">
                  <div
                    style="display: flex;flex-direction: column;align-items: center;justify-content: center;width: 73%;max-width: 92px;margin: auto;"
                  >
                    <div>
                      <a style="color: #6bb8ff" @click="drillDown(scope.row.ui[0])">&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;{{
                        scope.row.value[0]
                      }}&nbsp;&nbsp;</a>
                      <add-user-group :uid="scope.row.ui[0]" />
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column
                v-for="(v,k,index) in tableHeaderShow.slice(2,tableHeaderShow.length)"
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
                      <a style="color: #6bb8ff" @click="drillDown( scope.row.ui.slice(1,scope.row.ui.length)[k] )">
                        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
                        {{ scope.row.value.slice(1, scope.row.value.length)[k] }}
                        &nbsp;&nbsp;
                      </a>
                      <add-user-group :uid="scope.row.ui.slice(1,scope.row.ui.length)[k]" />
                    </div>

                    <div style="text-align: center;">
                      <template> {{
                        scope.row[lookTypScaleMap[lookTyp]].slice(1, scope.row[lookTypScaleMap[lookTyp]].length)[k]
                      }}%&nbsp;&nbsp;
                      </template>
                    </div>
                  </div>
                </template>
              </el-table-column>
            </page-table>
          </div>
          <div v-else-if="chartType==2">
            <retention-line2 v-if="g2Show" :is-scale="isScale" :chart-data="tableData" />
          </div>
          <div v-else-if="chartType==3">
            <retention-line v-if="g2Show" :is-scale="isScale" :chart-data="tableData" :x-data="getXdata" />
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
    'AddUserGroup': () => import('@/views/behavior-analysis/components/AddUserGroup'),
    'RetentionLine': () => import('@/components/Charts/RetentionLine'),
    'RetentionLine2': () => import('@/components/Charts/RetentionLine2'),
    'PageTable': () => import('@/components/PageTable'),
    'Date': () => import('@/components/AnalyseTools/FilterDate/Date')
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
    retentionRes: {
      type: Array,
      default: []
    },
    windowTime: {
      type: Number,
      default: 1
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
          value: '天',
          label: '天'
        }
      ],
      lookTypOpt: [
        {
          value: 'retention',
          label: '留存'
        },
        {
          value: 'lose',
          label: '流失'
        }
      ],
      form: {
        windowTime: this.windowTime,
        windowTimeFormat: '天'
      },
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
      tableData: []
    }
  },
  computed: {

    getXdata() {
      const xData = ['当天']
      for (let i = 0; i < this.form.windowTime; i++) {
        xData.push(`第${i + 1}天`)
      }

      return xData
    }
  },
  watch: {
    isScale: {
      deep: true,
      handler() {
        this.refreshData()
      }
    },
    'form.windowTime': {
      deep: true,
      handler() {
        this.$emit('changeWindowTime', this.form.windowTime, this.form.windowTimeFormat)
      }
    },
    'form.windowTimeFormat': {
      deep: true,
      handler() {
        this.$emit('changeWindowTime', this.form.windowTime, this.form.windowTimeFormat)
      }
    },
    retentionRes: {
      deep: true,
      handler() {
        this.init()
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
  beforeMount() {
    this.init()
  },
  methods: {
    drillDown(ui) {
      this.$store.dispatch('baseData/SETUI', ui)

      this.$router.push({ path: '/user-analysis/user_list' })
    },
    download(fName) {
      elTable2Excel(this, 'pagetable', `留存分析:${fName}`)
    },
    getTableHeader() {
      const header = ['日期', '初始事件触发用户数', '当天']
      for (let i = 0; i < this.form.windowTime; i++) {
        header.push(`第${i + 1}天`)
      }
      this.tableHeaderShow = header
    },

    changeWindow(input) {

    },
    changeWindowTime(input) {

    },

    NaN2Zero(v) {
      if (isNaN(v)) {
        return 0
      }
      return v
    },
    changeTable() {

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
      this.tableData = []
      for (const v of this.retentionRes) {
        const conversionScaleArr = []
        const washScaleArr = []
        const firstDayUserNum = v['value'][0]
        for (const k in v['value']) {
          if (k == 0) {
            conversionScaleArr.push(0)
            washScaleArr.push(0)
            continue
          }
          const conversionScale = this.conversionScale(v['value'][k], firstDayUserNum)
          conversionScaleArr.push(conversionScale)
          washScaleArr.push(this.washScale(conversionScale))
        }
        v['washScaleArr'] = washScaleArr
        v['conversionScaleArr'] = conversionScaleArr
        this.tableData.push(v)
      }
    },
    washScale(v1) {
      const tmp = 100 - parseFloat(v1)
      if (isNaN(tmp)) {
        return 0
      }
      return (tmp).toFixed(2)
    },
    conversionScale(v1, v2) {
      const tmp = (Number(v1 / v2) * 100)
      if (isNaN(tmp)) {
        return 0
      }
      return tmp.toFixed(2)
    },
    init() {
      this.getTableHeader()
      this.initTableData()
      this.refreshData()
    }
  }

}
</script>

<style scoped src="@/styles/retention-res.css"/>
