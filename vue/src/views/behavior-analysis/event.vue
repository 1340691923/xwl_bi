<template>
  <div style="display:flex;justify-content:space-between">
    <div class="content_xwl">
      <div class="header_xwl" style="background: white">
        <div class="root_xwl">
          <div class="main_xwl">
            <a-tooltip placement="right" style="cursor: pointer">
              <template slot="title">
                <span>分析某段时间内，某个事件或事件属性的整体趋势情况</span>
              </template>
              <span class="title_xwl" style="color: #202d3f">&nbsp;&nbsp;事件分析 <a-icon
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
            <report-table-list :rt_type="Number(1)" />
          </div>
        </div>
      </div>
      <split-pane :min-percent="0" :default-percent="28" split="vertical" @resize="onResize">
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
                  分析指标
                </div>

                <div class="xwl_main">

                  <draggable v-model="form.zhibiaoArr" handle=".drageTag" chosen-class="chosen" animation="1000">
                    <div v-for="(v,index) in form.zhibiaoArr" v-drag>
                      <div v-if="v.typ == 1" class="row___xwl" style="padding: 10px;">

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
                                  size="small"
                                  placeholder="请输入..."
                                  autofocus="autofocus"
                                  allow-clear
                                />
                              </a-tooltip>
                            </el-row>
                            <el-row style="padding-top: 10px" :span="6">
                              <a-select
                                v-model="form.zhibiaoArr[index].eventName"
                                size="default"
                                dropdown-match-select-width
                                show-search
                                default-active-first-option
                                style="width: 40%;"
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

                              <el-tag type="warning">的</el-tag>
                              <count-select
                                v-model="form.zhibiaoArr[index].selectAttr"
                                :options="form.zhibiaoArr[index].attrOptions"
                                placeholder="请筛选维度"
                              />
                            </el-row>
                          </el-col>
                          <el-col :span="4">
                            <a-tooltip placement="top" style="cursor: pointer">
                              <template slot="title">
                                <span>复制指标</span>
                              </template>
                              <a-button
                                type="link"
                                icon="copy"
                                style="margin-left: 4px"
                                class="right_icon_zhibiao"
                                @click="copyZhibiao(index)"
                              />
                            </a-tooltip>
                            <a-tooltip placement="top" style="cursor: pointer">
                              <template slot="title">
                                <span>删除指标</span>
                              </template>
                              <a-button
                                type="link"
                                icon="delete"
                                style="margin-left: 4px"
                                class="right_icon_zhibiao"
                                @click.native="form.zhibiaoArr.splice(index,1)"
                              />
                            </a-tooltip>
                          </el-col>
                        </el-row>

                        <div class="filters_xwl">
                          <filter-where
                            v-model="form.zhibiaoArr[index].relation"
                            font-color="#ffba00"
                            table-typ="2"
                            :data-type-map="attrMap"
                            :options="eventAttrOptions"
                          />
                        </div>

                      </div>
                      <div v-if="v.typ == 2" class="row___xwl" style="padding: 10px;">

                        <el-row>
                          <el-col :span="2">
                            <el-tag class="drageTag">{{ index + 1 }}</el-tag>
                          </el-col>
                          <el-col :span="18">
                            <el-row :span="15" class="zhibiao">
                              <el-col :span="13">
                                <a-tooltip placement="topLeft" style="cursor: pointer">
                                  <template slot="title">
                                    <span>点击修改公式名称</span>
                                  </template>
                                  <a-input
                                    :ref="('getFocus'+index).toString()"
                                    v-model:value="form.zhibiaoArr[index].eventNameDisplay"
                                    class="eventNameDisplayInput"
                                    size="small"
                                    placeholder="请输入..."
                                    autofocus="autofocus"
                                    allow-clear
                                    @change="getFocus1('getFocus'+index)"
                                  />
                                </a-tooltip>
                              </el-col>

                              <el-col :span="1">

                                <a-select
                                  v-model="form.zhibiaoArr[index].scaleType"
                                  size="small"
                                  default-value="1"
                                  style="width: 120px"
                                >
                                  <a-select-option value="1">
                                    两位小数
                                  </a-select-option>
                                  <a-select-option value="2">
                                    百分比
                                  </a-select-option>
                                  <a-select-option value="3">
                                    取整
                                  </a-select-option>
                                </a-select>
                              </el-col>
                            </el-row>
                            <el-row style="padding-top: 10px" :span="6">
                              <a-select
                                v-model="form.zhibiaoArr[index]['one'].eventName"
                                size="default"
                                dropdown-match-select-width
                                show-search
                                default-active-first-option
                                style="width: 120px"
                                @change="changeEventNameDisplay(index,'one')"
                              >

                                <a-select-option
                                  v-for="(v,k,index) in metaEventList"
                                  :key="index"
                                  :value="v.event_name"
                                >
                                  {{ v.show_name == '' ? v.event_name : v.show_name }}
                                </a-select-option>
                              </a-select>

                              <el-popover
                                placement="right"
                                trigger="click"
                              >
                                <filter-where
                                  v-model="form.zhibiaoArr[index]['one'].relation"
                                  table-typ="2"
                                  :data-type-map="attrMap"
                                  :options="eventAttrOptions"
                                />
                                <el-button slot="reference" type="text" icon="el-icon-data-board" />
                              </el-popover>

                              <el-tag>的</el-tag>

                              <count-select
                                v-model="form.zhibiaoArr[index]['one'].selectAttr"
                                :options="form.zhibiaoArr[index]['one'].attrOptions"
                                placeholder="请筛选维度"
                              />

                              <a-select v-model="form.zhibiaoArr[index].operate" size="small" style="width: 60px">
                                <a-select-option value="plus">
                                  加
                                </a-select-option>
                                <a-select-option value="minus">
                                  减
                                </a-select-option>
                                <a-select-option value="multiply">
                                  乘
                                </a-select-option>
                                <a-select-option value="divide">
                                  除
                                </a-select-option>
                              </a-select>

                              <a-select
                                v-model="form.zhibiaoArr[index]['two'].eventName"
                                size="default"
                                dropdown-match-select-width
                                show-search
                                default-active-first-option
                                style="width: 120px"
                              >
                                <a-select-option
                                  v-for="(v,k,index) in metaEventList"
                                  :key="index"
                                  :value="v.event_name"
                                >
                                  {{ v.show_name == '' ? v.event_name : v.show_name }}
                                </a-select-option>
                              </a-select>

                              <el-popover
                                placement="right"
                                trigger="click"
                              >
                                <filter-where

                                  v-model="form.zhibiaoArr[index]['two'].relation"
                                  table-typ="2"
                                  :data-type-map="attrMap"
                                  :options="eventAttrOptions"
                                />
                                <el-button slot="reference" type="text" icon="el-icon-data-board" />
                              </el-popover>

                              <el-tag>的</el-tag>
                              <count-select
                                v-model="form.zhibiaoArr[index]['two'].selectAttr"
                                :options="form.zhibiaoArr[index]['two'].attrOptions"
                                placeholder="请筛选维度"
                              />
                            </el-row>
                          </el-col>
                          <el-col :span="4">
                            <el-tooltip
                              :content="'被除数不分组: ' + form.zhibiaoArr[index].divisor_no_grouping"
                              placement="top"
                            >
                              <el-switch
                                v-model="form.zhibiaoArr[index].divisor_no_grouping"
                                active-color="#13ce66"
                                inactive-color="#ff4949"
                                :active-value="true"
                                :inactive-value="false"
                              />
                            </el-tooltip>
                            <a-tooltip placement="top" style="cursor: pointer">
                              <template slot="title">
                                <span>复制公式</span>
                              </template>
                              <a-button
                                type="link"
                                icon="copy"
                                style="margin-left: 4px;"
                                class="right_icon"
                                @click="copyZhibiao(index)"
                              />
                            </a-tooltip>
                            <a-tooltip placement="top" style="cursor: pointer">
                              <template slot="title">
                                <span>删除公式</span>
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
                      </div>
                    </div>
                  </draggable>
                </div>

                <div v-if="form.zhibiaoArr.length < 30" style="margin-top: 4px; padding: 0 12px;">
                  <span class="footadd___2D4YB" @click="addZhibiao(1)">
                    <a-icon type="plus" />
                    添加指标
                  </span>
                  <span class="footadd___2D4YB" @click="addZhibiao(2)">
                    <a-icon type="line-height" />
                    添加公式
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

              <div style="width: 100%;   padding-bottom: 8px;border-bottom: 1px solid #f0f2f5">
                <filter-group v-model="form.groupBy" limit="20" :options="eventAttrOptions" />
              </div>
            </div>

            <div
              style="width: 100%;height:  50px;margin-bottom: 0px;z-index: 10000;border-top: 1px  solid #f0f2f5;background: white;display: flex;align-items: center;justify-content: center"
            >
              <add-report-table
                :rt-type="Number(1)"
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
              <event-result
                v-if="eventResShow"
                ref="eventRes"
                v-model="form.date"
                :window-time-format="form.windowTimeFormat"
                style="padding: 20px"
                :event-res="eventRes"
                @changeWindowTime="changeWindowTime"
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

import { EventList, GetConfigs, LoadPropQuotas } from '@/api/analysis'
import { FindRtById } from '@/api/pannel'

export default {
  name: 'Event',
  components: {
    'CountSelect': () => import('@/components/AnalyseTools/CountSelect'),
    'FilterWhere': () => import('@/components/AnalyseTools/FilterWhere/index'),
    'FilterGroup': () => import('@/components/AnalyseTools/FilterGroup/index'),
    'ReportTableList': () => import('@/views/behavior-analysis/components/ReportTableList'),
    'AddReportTable': () => import('@/views/behavior-analysis/components/AddReportTable'),
    'EventResult': () => import('@/views/behavior-analysis/components/EventResult'),
    'FilterUserGroup': () => import('@/components/AnalyseTools/FilterUserGroup'),
    SmallSelect: () => import('@/components/AnalyseTools/FilterWhere/SmallSelect'),
    draggable
  },
  data() {
    return {
      reportTableName: '',
      currentReportTable: {
        name: '',
        remark: ''
      },
      spinning: false,
      drawerShow: false,

      eventRes: {},
      eventResShow: true,
      metaEventList: [],
      userAttr: [],
      eventAttr: [],
      form: {
        zhibiaoArr: [],
        groupBy: [],
        userGroup: [],
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
        date: [
          moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD'),
          moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD')
        ],
        windowTimeFormat: '按天'
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
    this.init()
  },

  methods: {
    download() {
      this.$refs['eventRes'].download('全量数据')
    },
    refreshRes() {
      this.$nextTick(() => {
        this.eventResShow = true
      })
    },
    changeWindowTime(windowTime) {
      this.form.windowTimeFormat = windowTime
      this.go()
    },
    onResize() {
      this.eventResShow = false
      this.debounceHandleSizeChange()
    },
    async getLoadPropQuotas(event_name) {
      const res = await LoadPropQuotas({ event_name: event_name, appid: this.$store.state.baseData.EsConnectID })

      const attrOptions = [
        {
          value: '默认',
          label: '默认',
          children: [
            {
              value: 'A1',
              label: '总次数'
            },
            {
              value: 'A2',
              label: '触发用户数'
            },
            {
              value: 'A3',
              label: '人均次数'
            }
          ]
        }
      ]
      if (res.code == 0) {
        for (const data of res.data) {
          const obj = {
            value: data.attribute_name,
            label: data.show_name,
            children: []
          }

          for (const k in data.analysis) {
            obj.children.push({
              value: k.toString(),
              label: data.analysis[k]
            })
          }
          attrOptions.push(obj)
        }
      }
      return attrOptions
    },
    async changeEventNameDisplay(index, parmas) {
      let eventNameDisplay = ''
      for (const v of this.metaEventList) {
        let eventName = ''
        if (parmas) {
          eventName = this.form.zhibiaoArr[index][parmas].eventName
        } else {
          eventName = this.form.zhibiaoArr[index].eventName
        }

        if (v.event_name == eventName) {
          eventNameDisplay = v.show_name == '' ? v.event_name : v.show_name
        }
      }

      const attrOptions = await this.getLoadPropQuotas(parmas ? this.form.zhibiaoArr[index][parmas].eventName : this.form.zhibiaoArr[index].eventName)

      this.form.zhibiaoArr[index].attrOptions = attrOptions
      this.form.zhibiaoArr[index].eventNameDisplay = eventNameDisplay
    },
    async initReportData() {
      const id = this.$route.params.id
      if (id != ':id' && Number(id) != 0) {
        const res = await FindRtById({ id: Number(id), appid: this.$store.state.baseData.EsConnectID })
        if (res.code != 0) {
          this.$message({
            offset: 60,
            type: 'error',
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
        this.addZhibiao(1)
      }
    },

    async getMetaEventList() {
      const res = await GetConfigs({ 'appid': this.$store.state.baseData.EsConnectID })

      if (res.code != 0) {
        this.$message({
          offset: 60,
          type: 'error',
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
    async addZhibiao(typ) {
      if (this.form.zhibiaoArr.length >= 30) return
      const attrOptions = await this.getLoadPropQuotas(this.metaEventList[0].event_name)

      if (this.metaEventList.length > 0) {
        switch (typ) {
          case 1:
            const obj = {
              'selectAttr': [],
              'attrOptions': [],
              'typ': 1,
              'eventName': this.metaEventList[0].event_name,
              'eventNameDisplay': this.metaEventList[0].show_name != '' ? this.metaEventList[0].show_name : this.metaEventList[0].event_name,
              'relation': {
                filterType: 'COMPOUND',
                filts: [],
                relation: '且'
              }
            }

            obj['attrOptions'] = attrOptions

            this.form.zhibiaoArr.push(obj)
            break
          case 2:
            const obj2 = {
              'typ': 2,
              'scaleType': '1',
              'operate': 'divide',
              'divisor_no_grouping': true,
              'eventNameDisplay': this.metaEventList[0].show_name != '' ? this.metaEventList[0].show_name : this.metaEventList[0].event_name,
              'one': {
                'selectAttr': [],
                'attrOptions': attrOptions,
                'eventName': this.metaEventList[0].event_name,
                'relation': {
                  filterType: 'COMPOUND',
                  filts: [],
                  relation: '且'
                }
              },
              'two': {
                'selectAttr': [],
                'attrOptions': attrOptions,
                'eventName': this.metaEventList[0].event_name,
                'relation': {
                  filterType: 'COMPOUND',
                  filts: [],
                  relation: '且'
                }
              }
            }

            this.form.zhibiaoArr.push(obj2)
        }
      }
    },
    copyZhibiao(index) {
      if (this.form.zhibiaoArr.length >= 30) return
      console.log('this.form.zhibiaoArr[index]', index)
      const obj = JSON.parse(JSON.stringify(this.form.zhibiaoArr[index]))
      this.form.zhibiaoArr.splice(index, 0, obj)
    },
    getFocus1(f) {
      this.$nextTick(() => {
        this.$refs[f][0].focus()
      })
    },
    async go() {
      this.spinning = true
      this.eventResShow = false
      const form = JSON.parse(JSON.stringify(this.form))

      for (const i in form.zhibiaoArr) {
        if (form.zhibiaoArr[i].typ == 1) {
          if (form.zhibiaoArr[i].attrOptions) {
            delete form.zhibiaoArr[i].attrOptions
          }
        } else if (form.zhibiaoArr[i].typ == 2) {
          if (form.zhibiaoArr[i]['one'].attrOptions) {
            delete form.zhibiaoArr[i]['one'].attrOptions
          }
          if (form.zhibiaoArr[i]['two'].attrOptions) {
            delete form.zhibiaoArr[i]['two'].attrOptions
          }
        }
      }

      form['appid'] = this.$store.state.baseData.EsConnectID
      const res = await EventList(form)
      if (res.code != 0) {
        this.$message({
          offset: 60,
          type: 'error',
          message: res.msg
        })
        this.eventRes = {}
        this.spinning = false
        this.$nextTick(() => {
          this.eventResShow = true
        })
        return
      }

      this.eventRes = res.data

      this.spinning = false
      this.$nextTick(() => {
        this.eventResShow = true
      })
    }
  }
}
</script>

<style scoped src="@/styles/event.css"/>

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
