<template>
  <div style="display:flex;justify-content:space-between">
    <div class="content_xwl">
      <div class="header_xwl" style="background: white">
        <div class="root_xwl">
          <div class="main_xwl">
            <a-tooltip placement="right" style="cursor: pointer">
              <template slot="title">
                <span>以某段时间做过步骤1的用户为样本，查看窗口期内，指定步骤下用户的转化情况</span>
              </template>
              <span class="title_xwl" style="color: #202d3f">&nbsp;&nbsp;漏斗分析 <a-icon
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
            <report-table-list :rt_type="Number(3)" />
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
                  漏斗步骤
                </div>

                <div class="xwl_main">

                  <draggable v-model="form.zhibiaoArr" handle=".drageTag" chosen-class="chosen" animation="1000">

                    <div v-for="(v,index) in form.zhibiaoArr" v-drag>

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
                                v-model="form.zhibiaoArr[index].eventName"
                                dropdown-match-select-width
                                show-search
                                default-active-first-option
                                style="width: 75%;"
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
                          <el-col :span="4">
                            <a-tooltip placement="top" style="cursor: pointer">
                              <template slot="title">
                                <span>插入漏斗步骤</span>
                              </template>
                              <a-button
                                type="link"
                                icon="copy"
                                style="margin-left: 4px;"
                                class="right_icon"
                                @click="copyZhibiao(index)"
                              />
                            </a-tooltip>
                            <a-tooltip v-if="index>=2" placement="top" style="cursor: pointer">
                              <template slot="title">
                                <span>删除指标</span>
                              </template>
                              <a-button
                                type="link"
                                icon="delete"
                                style="margin-left: 4px;"
                                class="right_icon"
                                @click.native="form.zhibiaoArr.splice(index,1)"
                              />
                            </a-tooltip>
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
                  </draggable>
                </div>

                <div v-if="form.zhibiaoArr.length < 30" style="margin-top: 4px; padding: 0 12px;">
                  <span class="footadd___2D4YB" @click="addZhibiao()">
                    <a-icon type="plus" />
                    添加步骤指标
                  </span>
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
                <filter-user-group v-model="form.userGroup" />
              </div>

              <div v-show="true" style="width: 100%;   padding-bottom: 8px;border-bottom: 1px solid #f0f2f5">
                <filter-group v-model="form.groupBy" :options="eventAttrOptions" />
              </div>
              <div
                style="line-height: 18px;font-weight: 500; font-size: 13px; padding: 10px 16px 12px;font-weight: bolder"
              >
                分析窗口期
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

            <div
              style="width: 100%;height:  50px;margin-bottom: 0px;z-index: 10000;border-top: 1px  solid #f0f2f5;background: white;display: flex;align-items: center;justify-content: center"
            >
              <add-report-table
                :rt-type="Number(3)"
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
              <funnel-result
                v-if="funnelResShow"
                ref="funnelRes"
                v-model="form.date"
                style="padding: 20px"
                :padding="'20px'.toString()"
                :table-header="tableHeader"
                :group-data="groupData"
                :funnel-res="funnelRes"
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
import draggable from 'vuedraggable'
import moment from 'moment'

import { FunnelList, GetConfigs } from '@/api/analysis'
import { FindRtById } from '@/api/pannel'

export default {
  name: 'Funnel',
  components: {
    'FilterWhere': () => import('@/components/AnalyseTools/FilterWhere/index'),
    'FilterGroup': () => import('@/components/AnalyseTools/FilterGroup/index'),
    'FunnelResult': () => import('@/views/behavior-analysis/components/FunnelResult'),
    'ReportTableList': () => import('@/views/behavior-analysis/components/ReportTableList'),
    'AddReportTable': () => import('@/views/behavior-analysis/components/AddReportTable'),
    'FilterUserGroup': () => import('@/components/AnalyseTools/FilterUserGroup'),
    draggable
  },
  data() {
    return {

      currentReportTable: {
        name: '',
        remark: ''
      },
      spinning: false,
      drawerShow: false,
      tableHeader: [],
      groupData: {},
      funnelRes: [],
      funnelResShow: true,
      prevCount: 0,
      metaEventList: [],
      userAttr: [],
      eventAttr: [],
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
      reportTableName: '',
      form: {
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
      this.$refs['funnelRes'].download('全量数据')
    },
    refreshRes() {
      this.$nextTick(() => {
        this.funnelResShow = true
      })
    },
    onResize() {
      this.funnelResShow = false
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
        this.go()
      }
    },

    async init() {
      await this.getMetaEventList()
      if (this.form.zhibiaoArr.length == 0) {
        this.addZhibiao()
        this.addZhibiao()
      }
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
      this.metaEventList = res.data.event_name_list

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
    copyZhibiao(index) {
      if (this.form.zhibiaoArr.length >= 30) return
      const obj = JSON.parse(JSON.stringify(this.form.zhibiaoArr[index]))
      this.form.zhibiaoArr.splice(index, 0, obj)
    },
    getFocus1(f) {
      this.$nextTick(function() {
        this.$refs[f][0].focus()
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
                obj.ui = v.ui
                obj.count = v.count
                obj.succScale = 100.00
                obj.lostScale = 0
              } else {
                obj.ui = v.ui
                obj.count = v.count
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
    getGroupByLable(v) {
      for (const option of this.allAttrOptions[0].options) {
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

<style scoped src="@/styles/funnel.css"/>

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
