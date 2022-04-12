<template>
  <div style="display:flex;justify-content:space-between">
    <div class="content_xwl">
      <div class="header_xwl" style="background: white">
        <div class="root_xwl">
          <div class="main_xwl">
            <a-tooltip placement="right" style="cursor: pointer">
              <template slot="title">
                <span>以一个事件为终点，通过桑基图直观掌握用户行为扩展路线</span>
              </template>
              <span class="title_xwl" style="color: #202d3f">&nbsp;&nbsp;智能路径分析 <a-icon
                type="question-circle"
              />
                <template v-if="reportTableName!=''">
                  {{ reportTableName }}
                </template>
              </span>
            </a-tooltip>
          </div>
          <div class="actions_xwl">
            <a-tooltip placement="top" style="cursor: pointer">
              <template slot="title">
                <span>以页面格式下载全量数据</span>
              </template>
              <a-button type="link" class="actions_xwl_btn" icon="download" @click="download" />
            </a-tooltip>
            <report-table-list :rt_type="Number(4)" />
          </div>
        </div>
      </div>
      <split-pane :min-percent="0" :default-percent="22" split="vertical" @resize="onResize">
        <template slot="paneL">
          <div
            id="scollL"
            style="height: 95%;width: 100px;display: inline-block; height: 100%;vertical-align: top;width: 100%;background: white;"
          >
            <div style="width: 100%;height: calc(100% - 140px); overflow-x: hidden; overflow-y: auto;">

              <div style="width: 100%;   padding-bottom: 8px;border-bottom: 1px solid #f0f2f5">
                <div
                  style="line-height: 18px;font-weight: 500; font-size: 13px; padding: 10px 16px 12px;font-weight: bolder"
                >
                  参与分析的事件
                </div>
                <div class="xwl_main">
                  <div>
                    <div class="row___xwl" style="padding: 10px;">
                      <select2 v-if="eventNameSelectShow" v-model="form.eventNames" :options="eventOptions" :checkeds="false" placeholder="请选择事件名" />
                    </div>
                  </div>
                </div>
              </div>

              <div style="width: 100%;   padding-bottom: 8px;border-bottom: 1px solid #f0f2f5">
                <div
                  style="line-height: 18px;font-weight: 500; font-size: 13px; padding: 10px 16px 12px;font-weight: bolder"
                >
                  终点事件
                </div>
                <div class="xwl_main">
                  <div v-for="(v,index) in form.zhibiaoArr" v-if="index == 0">
                    <div class="row___xwl" style="padding: 10px;">
                      <el-row>
                        <el-col :span="2">
                          <el-tag type="warning" class="drageTag">{{ index + 1 }}</el-tag>
                        </el-col>
                        <el-col :span="18">
                          <el-row :span="15" class="zhibiao">
                            <a-tooltip placement="topLeft" style="cursor: pointer">
                              <template slot="title">
                                <span>点击修改指标名称</span>
                              </template>
                              <a-input
                                :ref="('getFocus'+index).toString()"
                                v-model:value="form.zhibiaoArr[index].eventNameDisplay"
                                class="eventNameDisplayInput"
                                placeholder="请输入..."
                                autofocus="autofocus"
                                allow-clear
                                @change="getFocus1('getFocus'+index)"
                              />
                            </a-tooltip>
                          </el-row>
                          <el-row style="padding-top: 5px" :span="6">
                            <a-select
                              style="width: 75%;"
                              v-model="form.zhibiaoArr[index].eventName"
                              dropdown-match-select-width
                              show-search
                              default-active-first-option

                              @change="changeEventNameDisplay(index)"
                            >

                              <a-select-option
                                v-for="(v,k,index) in metaEventList"
                                :key="index"
                                :value="v.event_name"
                              >
                                {{ v.show_name == '' ? v.event_name : v.show_name }}
                              </a-select-option>

                            </a-select>
                          </el-row>
                        </el-col>
                      </el-row>
                      <div class="filters_xwl">
                        <filter-where
                          v-model="form.zhibiaoArr[index].relation"
                          table-typ="2"
                          :data-type-map="attrMap"
                          :options="eventAttrOptions"
                        />
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div style="width: 100%;   padding-bottom: 8px;border-bottom: 1px solid #f0f2f5">
                <div
                  style="line-height: 18px;font-weight: 500; font-size: 13px; padding: 10px 16px 10px;font-weight: bolder"
                >
                  全局筛选事件维度
                </div>

                <filter-where
                  v-model="form.whereFilter"
                  :data-type-map="attrMap"
                  table-typ="2"
                  :options="eventAttrOptions"
                />
              </div>

              <div style="width: 100%;   padding-bottom: 8px;border-bottom: 1px solid #f0f2f5">
                <div
                  style="line-height: 18px;font-weight: 500; font-size: 13px; padding: 10px 16px 10px;font-weight: bolder"
                >
                  全局筛选用户维度
                </div>
                <filter-where
                  v-model="form.whereFilterByUser"
                  :data-type-map="attrMap"
                  table-typ="1"
                  :options="userAttrOptions"
                />
              </div>
              <div style="width: 100%;   padding-bottom: 8px;border-bottom: 1px solid #f0f2f5">
                <div
                  style="line-height: 18px;font-weight: 500; font-size: 13px; padding: 10px 16px 12px;font-weight: bolder"
                >
                  初始事件与最终事件最大间隔时间
                </div>
                <div style="padding-left: 8px;margin-left: 10px">
                  <el-input-number
                    v-model="form.windowTime"
                    size="mini"
                    controls-position="right"
                    style="width: 120px"
                    :min="1"
                  />
                  <el-select v-model="form.windowTimeFormat" size="mini" style="width: 100px">
                    <el-option
                      v-for="item in windowTimeOpt"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    />
                  </el-select>
                </div>
              </div>

              <div style="width: 100%;   padding-bottom: 8px;border-bottom: 1px solid #f0f2f5">
                <filter-user-group v-model="form.userGroup" />
              </div>
            </div>

            <div
              style="width: 100%;height:  50px;margin-bottom: 0px;z-index: 10000;border-top: 1px  solid #f0f2f5;background: white;display: flex;align-items: center;justify-content: center"
            >
              <add-report-table
                :rt-type="Number(4)"
                :name="currentReportTable.name"
                :remark="currentReportTable.remark"
                :data="this.form"
                style="width: 200px;height: 50px;line-height: 50px;margin: 0px auto"
                @go="go"
              />
            </div>

          </div>
        </template>
        <template slot="paneR">
          <a-spin tip="计算中..." :spinning="spinning">
            <div class="spin-content">
              <trace-result
                v-if="traceResShow"
                ref="traceRes"
                v-model="form.date"
                style="padding: 20px"
                :trace-charts-res="traceChartsRes"
                :trace-table-res="traceTableRes"
                @go="go"
              />
            </div>
          </a-spin>
        </template>
      </split-pane>
    </div>
  </div>
</template>

<script>
import { debounce } from 'lodash'

import moment from 'moment'

import { GetConfigs, TraceList } from '@/api/analysis'
import { FindRtById } from '@/api/pannel'

export default {
  name: 'Trace',
  components: {
    'FilterWhere': () => import('@/components/AnalyseTools/FilterWhere/index'),
    'FilterGroup': () => import('@/components/AnalyseTools/FilterGroup/index'),
    'TraceResult': () => import('@/views/behavior-analysis/components/TraceResult'),
    'ReportTableList': () => import('@/views/behavior-analysis/components/ReportTableList'),
    'AddReportTable': () => import('@/views/behavior-analysis/components/AddReportTable'),
    'FilterUserGroup': () => import('@/components/AnalyseTools/FilterUserGroup'),
    'Select2': () => import('@/components/AnalyseTools/Select2/Select')
  },
  data() {
    return {
      eventNameSelectShow: true,
      eventOptions: [],
      traceTableRes: [],
      traceChartsRes: [],
      windowTimeOpt: [
        {
          value: '天',
          label: '天'
        }, {
          value: '小时',
          label: '小时'
        }, {
          value: '分钟',
          label: '分钟'
        }, {
          value: '秒',
          label: '秒'
        }
      ],
      currentReportTable: {
        name: '',
        remark: ''
      },
      spinning: false,
      drawerShow: false,
      tableHeader: [],
      traceRes: [],
      traceResShow: true,
      prevCount: 0,
      metaEventList: [],
      userAttr: [],
      eventAttr: [],
      reportTableName: '',
      form: {
        eventNames: [],
        zhibiaoArr: [],
        groupBy: [],

        whereFilter: {
          filterType: 'COMPOUND',
          filts: [],
          relation: '且'
        },

        whereFilterByUser: {
          filterType: 'COMPOUND',
          filts: [],
          relation: '且'
        },
        windowTime: 1,
        windowTimeFormat: '天',
        date: [
          moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD'),
          moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD')
        ]
      },

      eventAttrOptions: [],
      allAttrOptions: [],
      userAttrOptions: [],
      attrMap: [],
      debounceHandleSizeChange: undefined
    }
  },
  async beforeMount() {
    await this.initReportData()
    this.debounceHandleSizeChange = debounce(this.refreshRes, 500)
  },
  mounted() {
    this.init()
  },
  methods: {

    download() {
      this.$refs['traceRes'].download('全量数据')
    },
    refreshRes() {
      this.$nextTick(() => {
        this.traceResShow = true
      })
    },
    onResize() {
      this.traceResShow = false
      this.debounceHandleSizeChange()
    },
    changeEventNameDisplay(index) {
      let eventNameDisplay = ''
      for (const v of this.metaEventList) {
        if (v.event_name == this.form.zhibiaoArr[index].eventName) {
          eventNameDisplay = v.show_name == '' ? v.event_name : v.show_name
        }
      }
      this.form.zhibiaoArr[index].eventNameDisplay = eventNameDisplay
    },
    async initReportData() {
      const id = this.$route.params.id
      if (id != ':id' && Number(id) != 0) {
        const res = await FindRtById({ id: Number(id), 'appid': this.$store.state.baseData.EsConnectID })
        if (res.code != 0) {
          this.$message({
            type: 'error',
            offset: 60,
            message: res.msg
          })
          return
        }
        this.reportTableName = res.data.name
        this.currentReportTable.name = res.data.name
        this.currentReportTable.remark = res.data.remark
        this.form = JSON.parse(res.data.data)
        this.form.date = [
          moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD'),
          moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD')
        ]
        this.eventNameSelectShow = false
        this.$nextTick(() => {
          this.eventNameSelectShow = true
        })

        this.go()
      }
    },

    async init() {
      await this.getMetaEventList()
      if (this.form.zhibiaoArr.length == 0) {
        this.addZhibiao()
      }
    },

    async getMetaEventList() {
      this.eventOptions = []

      const res = await GetConfigs({ 'appid': this.$store.state.baseData.EsConnectID })

      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })

        return
      }
      this.metaEventList = res.data.event_name_list

      for (const i in this.metaEventList) {
        const obj = {}
        obj['label'] = this.metaEventList[i].show_name == '' ? this.metaEventList[i].event_name : this.metaEventList[i].show_name
        obj['value'] = this.metaEventList[i].event_name
        this.eventOptions.push(obj)
      }

      const attributeMap = res.data.attributeMap
      const allAttrOptions = []
      const eventAttrOptions = []
      const userAttrOptions = []
      this.attrMap = []
      this.attrMap = attributeMap
      const eventData = { label: '事件', options: [] }
      const userData = { label: '用户', options: [] }
      if (attributeMap.hasOwnProperty('2')) {
        for (const v of attributeMap['2']) {
          eventData.options.push({
            value: v.attribute_name,
            label: v.show_name == '' ? v.attribute_name : v.show_name
          })
        }
      }
      if (attributeMap.hasOwnProperty('1')) {
        for (const v of attributeMap['1']) {
          userData.options.push({
            value: v.attribute_name,
            label: v.show_name == '' ? v.attribute_name : v.show_name
          })
        }
      }

      eventAttrOptions.push(eventData)
      userAttrOptions.push(userData)
      allAttrOptions.push(eventData, userData)
      this.userAttrOptions = userAttrOptions
      this.eventAttrOptions = eventAttrOptions
      this.allAttrOptions = allAttrOptions
    },
    moment,
    addZhibiao() {
      if (this.form.zhibiaoArr.length >= 30) return

      this.form.zhibiaoArr.push({
        'eventName': this.metaEventList[0].event_name,
        'eventNameDisplay': this.metaEventList[0].show_name != '' ? this.metaEventList[0].show_name : this.metaEventList[0].event_name,
        'relation': {
          filterType: 'COMPOUND',
          filts: [],
          relation: '且'
        }
      })
    },
    getFocus1(f) {
      this.$nextTick(function() {
        this.$refs[f][0].focus()
      })
    },
    async go() {
      this.traceResShow = false
      this.spinning = true

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

<style scoped src="@/styles/trace.css"/>

<style>
.eventNameDisplayInput .ant-input {
  resize: none;
  border: none;
}

.eventNameDisplayInput .ant-input:focus {
  border: none;
  box-shadow: none;
}
</style>
