<template>
  <div style="display:flex;justify-content:space-between">
    <div class="content_xwl">
      <div class="header_xwl" style="background: white">
        <div class="root_xwl">
          <div class="main_xwl">
            <a-tooltip placement="right" style="cursor: pointer">
              <template slot="title">
                <span>分析当前状态下，指定用户的用户属性分布</span>
              </template>
              <span class="title_xwl" style="color: #202d3f">&nbsp;&nbsp;用户属性分析 <a-icon
                type="question-circle"
              />
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
            <report-table-list :rt_type="Number(5)" />
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
                  用户属性分析
                </div>

                <div class="xwl_main">
                  <count-select v-model="form.zhibiaoArr" :options="userAttrZhibiaoOptions" placeholder="请筛选用户属性" />
                </div>

              </div>

              <div style="width: 100%;   padding-bottom: 8px;border-bottom: 1px solid #f0f2f5">
                <div
                  style="line-height: 18px;font-weight: 500; font-size: 13px; padding: 10px 16px 10px;font-weight: bolder"
                >
                  且用户符合
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
                <filter-group v-model="form.groupBy" type-tag="用户" :options="userAttrOptions" />
              </div>

            </div>

            <div
              style="width: 100%;height:  50px;margin-bottom: 0px;z-index: 10000;border-top: 1px  solid #f0f2f5;background: white;display: flex;align-items: center;justify-content: center"
            >
              <add-report-table
                :rt-type="Number(5)"
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
              <user-attr-result
                v-if="userAttrResShow"
                ref="userAttrRes"
                :tablelabel="getTableLabel"
                style="padding: 20px"
                :user-attr-res="userAttrRes"
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

import { GetConfigs, UserAttrList } from '@/api/analysis'
import { FindRtById } from '@/api/pannel'

export default {
  name: 'Index',
  components: {
    'UserAttrResult': () => import('@/views/behavior-analysis/components/UserAttrResult'),
    'CountSelect': () => import('@/components/AnalyseTools/CountSelect'),
    'FilterWhere': () => import('@/components/AnalyseTools/FilterWhere/index'),
    'FilterGroup': () => import('@/components/AnalyseTools/FilterGroup/index'),
    'ReportTableList': () => import('@/views/behavior-analysis/components/ReportTableList'),
    'AddReportTable': () => import('@/views/behavior-analysis/components/AddReportTable'),
    'FilterUserGroup': () => import('@/components/AnalyseTools/FilterUserGroup'),
    draggable
  },
  data() {
    return {
      userAttrZhibiaoOptions: [],
      currentReportTable: {
        name: '',
        remark: ''
      },
      spinning: false,
      tableHeader: [],
      groupData: {},
      userAttrRes: [],
      userAttrResShow: true,
      metaEventList: [],
      userAttr: [],
      eventAttr: [],
      form: {
        userGroup: [],
        zhibiaoArr: [],
        groupBy: [],
        whereFilterByUser: {
          filterType: 'COMPOUND',
          filts: [],
          relation: '且'
        }
      },
      eventAttrOptions: [],
      allAttrOptions: [],
      userAttrOptions: [],
      attrMap: [],
      debounceHandleSizeChange: undefined
    }
  },
  computed: {
    getTableLabel() {
      let tmp = ''
      switch (this.form.zhibiaoArr[1]) {
        case '1':
          tmp = '用户数'
          break
        case '2':
          tmp = '总和'
          break
        case '3':
          tmp = '均值'
          break
        case '4':
          tmp = '人均值'
          break
        case '5':
          tmp = '中位值'
          break
        case '6':
          tmp = '最大值'
          break
        case '7':
          tmp = '最小值'
        case '8':
          tmp = '去重值'
          break
      }
      return this.form.zhibiaoArr[0] + '.' + tmp
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
    async initReportData() {
      const id = this.$route.params.id
      if (id != ':id' && Number(id) != 0) {
        const res = await FindRtById({ id: Number(id), 'appid': this.$store.state.baseData.EsConnectID })
        if (res.code != 0) {
          this.$message({
            offset: 60,
            type: 'error',
            message: res.msg
          })
          return
        }
        this.currentReportTable.name = res.data.name
        this.currentReportTable.remark = res.data.remark
        this.form = JSON.parse(res.data.data)
      }
    },
    download() {
      this.$refs['userAttrRes'].download('全量数据')
    },
    refreshRes() {
      this.$nextTick(() => {
        this.userAttrResShow = true
      })
    },
    onResize() {
      this.userAttrResShow = false
      this.debounceHandleSizeChange()
    },

    async init() {
      await this.getMetaEventList()
      if (this.form.zhibiaoArr.length == 0) {
        this.addZhibiao()
      }
    },

    async getMetaEventList() {
      this.userAttrZhibiaoOptions = [
        {
          value: '默认',
          label: '默认',
          children: [
            {
              value: '1',
              label: '用户数'
            }
          ]
        }
      ]

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

      if (attributeMap.hasOwnProperty('1')) {
        for (const v of attributeMap['1']) {
          userData.options.push({
            value: v.attribute_name,
            label: v.show_name == '' ? v.attribute_name : v.show_name
          })

          switch (v.data_type) {
            case 3:
              this.userAttrZhibiaoOptions.push({
                value: v.attribute_name,
                label: v.show_name == '' ? v.attribute_name : v.show_name,
                children: [
                  {
                    value: '8',
                    label: '去重值'
                  }
                ]
              })
              break
            case 1:
            case 2:
              this.userAttrZhibiaoOptions.push({
                value: v.attribute_name,
                label: v.show_name == '' ? v.attribute_name : v.show_name,
                children: [
                  {
                    value: '2',
                    label: '总和'
                  },
                  {
                    value: '3',
                    label: '均值'
                  },
                  {
                    value: '4',
                    label: '人均值'
                  },
                  {
                    value: '5',
                    label: '中位值'
                  },
                  {
                    value: '6',
                    label: '最大值'
                  },
                  {
                    value: '7',
                    label: '最小值'
                  },
                  {
                    value: '8',
                    label: '去重值'
                  }
                ]
              })
              break
          }
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
      if (this.form.zhibiaoArr.length >= 1) return
      this.form.zhibiaoArr = ['默认', '1']
    },
    getFocus1(f) {
      this.$nextTick(function() {
        this.$refs[f][0].focus()
      })
    },

    getGroupByLable(v) {
      for (const option of this.allAttrOptions[1].options) {
        if (option.value == v) {
          return option.label
        }
      }
    },
    async go() {
      this.spinning = true
      this.userAttrResShow = false
      const form = this.form
      form['appid'] = this.$store.state.baseData.EsConnectID
      const res = await UserAttrList(form)
      if (res.code != 0) {
        this.$message({
          offset: 60,
          type: 'error',
          message: res.msg
        })
        this.userAttrRes = []
        this.groupData = []
        this.spinning = false
        this.$nextTick(() => {
          this.userAttrResShow = true
        })
        return
      }
      this.userAttrRes = res.data.tableRes
      this.spinning = false
      this.$nextTick(() => {
        this.userAttrResShow = true
      })
    }
  }

}
</script>

<style scoped src="@/styles/user-as.css"/>
