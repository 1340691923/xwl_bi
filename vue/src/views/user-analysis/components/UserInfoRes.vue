<template>
  <div class="right_res">
    <div>
      <div id="top" />
      <div
        class="app-container"
        style="height: 100%; padding: 30px; background: white;"
      >
        <div style="display: flex; align-items: center; justify-content: space-between;">
          <div class="echartBox_title">
            <date v-model="form.date" @changeDate="filterDateCall" />

            <el-select v-model="form.windowTimeFormat" size="mini" style="width: 100px">

              <el-option
                v-for="item in windowTimeOpt"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
            <select2 v-model="form.eventNames" :options="eventOptions" :checkeds="true" placeholder="请选择事件名" />

            <el-button type="success" size="mini" @click="go">GO-></el-button>
          </div>
        </div>
        <el-row :gutter="20" style="margin-top: 15px">
          <el-col :span="14">
            <event-count-bar
              v-if="pageshow"
              :chart-data="chartData.barData"
              class-name="event-count-bar"
            />
          </el-col>
          <el-col :span="10">
            <event-count-pie
              v-if="pageshow"
              :chart-data="chartData.pieData"
              class-name="event-count-pie"
            />
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="16">
            <el-button
              v-if="infoForm.sortRule == 'desc'"
              size="mini"
              type="primary"
              icon="el-icon-sort-up"
              @click="changeOrderBy('asc')"
            >时间升序
            </el-button>
            <el-button
              v-if="infoForm.sortRule == 'asc'"
              size="mini"
              type="success"
              icon="el-icon-sort-down"
              @click="changeOrderBy('desc')"
            >时间降序
            </el-button>
            <el-select v-model="infoForm.eventName" size="mini" clearable placeholder="锁定事件" @change="changeLockEvent">
              <el-option v-for="(item, index) in eventOptions" :key="index" :label="item.label" :value="item.value">
                <span style="float: left">{{ item.label }}</span>
                <span style="float: right; color: #8492a6; font-size: 13px">{{ item.value }}</span>
              </el-option>
            </el-select>
          </el-col>
          <el-col :span="8" />
        </el-row>
        <div style="margin-top: 30px">
          <template v-if=" Object.keys(userEventTimeList).length > 0 ">
            <div v-for="(v,year,index) in userEventTimeList" v-if="dateListShow" :key="index">
              <el-divider content-position="left"><i class="el-icon-date" />{{ year }}</el-divider>
              <el-timeline>
                <el-timeline-item
                  v-for="(v2,k2,index2) in v"
                  :key="index2"
                  color="#0bbd87"
                  :timestamp="v2.date_t"
                  placement="top"
                >
                  <el-collapse>
                    <el-collapse-item :name="index2">
                      <template slot="title">{{ v2.xwl_part_event_desc }}<i class="header-icon el-icon-thumb" />
                      </template>
                      <el-descriptions>
                        <el-descriptions-item
                          v-for="(v4,k4,index4) in v2"
                          v-if="!inBanAttr(k4)"
                          :key="index4"
                          :label="eventAttrMap[k4] == undefined || eventAttrMap[k4] == '' ?k4:eventAttrMap[k4]"
                        >
                          {{ v4 }}
                        </el-descriptions-item>
                      </el-descriptions>
                    </el-collapse-item>
                  </el-collapse>
                </el-timeline-item>
              </el-timeline>
            </div>

          </template>
          <a-empty v-else>
            <span slot="description">请重新筛选条件</span>
          </a-empty>
        </div>
        <el-divider style="margin-bottom: 60px;">
          <el-button v-loading="loadingMore" icon="el-icon-arrow-down" type="primary" @click="more">加载更多</el-button>
        </el-divider>

        <el-backtop target=".right_res" :bottom="100" :right="100">
          <div
            style="{
        height: 100%;
        width: 100%;
        background-color: #f2f5f6;
        box-shadow: 0 0 6px rgba(0,0,0, .12);
        text-align: center;
        line-height: 40px;
        color: #1989fa;
      }"
          >
            UP
          </div>
        </el-backtop>
      </div>

    </div>

  </div>
</template>

<script>
import { GetConfigs, UserEventCountList, UserEventDetailList } from '@/api/analysis'
import moment from 'moment'

export default {
  name: 'UserInfoRes',
  components: {
    'EventCountBar': () => import('@/components/Charts/EventCountBar'),
    'EventCountPie': () => import('@/components/Charts/EventCountPie'),
    'Select2': () => import('@/components/AnalyseTools/Select2/Select'),
    'Date': () => import('@/components/AnalyseTools/FilterDate/Date')
  },
  props: {
    userId: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      loadingMore: false,
      dateListShow: true,
      userEventTimeList: {},
      chartData: {
        pieData: [],
        barData: []
      },
      form: {
        date: [],
        eventNames: [],
        userId: this.userId,
        windowTimeFormat: '按天'
      },
      eventName: '',
      infoForm: {
        sortRule: 'desc',
        date: [],
        eventName: '',
        userId: this.userId,
        windowTimeFormat: '按天'
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
        }
      ],
      eventOptions: [],
      pageshow: false,
      page: 1,
      page_size: 30,
      eventAttrMap: {}
    }
  },
  computed: {},
  watch: {},

  beforeMount() {
    this.init()
  },
  methods: {
    refreshDateListShow() {
      this.dateListShow = false
      this.$nextTick(() => {
        this.dateListShow = true
      })
    },
    inBanAttr(attr) {
      const banAttrList = ['xwl_part_event', 'xwl_part_date', 'date_year', 'date_t', 'xwl_part_event_desc']
      return banAttrList.indexOf(attr) !== -1
    },
    refreshPage() {
      this.pageshow = false
      this.$nextTick(() => {
        this.pageshow = true
      })
    },
    async more() {
      this.loadingMore = true
      if(this.page == 1){
        this.page = this.page + 1
      }
      const status = await this.getDescList()
      this.loadingMore = false
      if (status == 1) {
        this.page = this.page + 1
      }
    },
    async getDescList(lockEvent) {
      if (this.page == 1) {
        this.userEventTimeList = {}
      }
      const form = {
        eventNames: this.form.eventNames,
        page: this.page,
        page_size: this.page_size,
        appid: this.$store.state.baseData.EsConnectID,
        orderBy: this.infoForm.sortRule,
        eventName: this.infoForm.eventName,
        userId: this.userId,
        date: this.form.date
      }

      const res = await UserEventDetailList(form)

      if (res.code != 0) {
        this.$message({
          offset: 60,
          type: 'error',
          message: res.msg
        })
        this.userEventTimeList = {}
        this.refreshDateListShow()
        return new Promise(resolve => {
          resolve(0)
        })
      }

      if(res.data.list != null){
        let tmp = {}
        for(let v of res.data.list){
          let key = Object.keys(v)[0]
          tmp[key] = v[key]
        }
        res.data.list = tmp
      }

      if (res.data.list == null || Object.keys(res.data.list).length == 0) {
        res.data.list = {}
        this.$message({
          offset: 60,
          type: 'error',
          message: '暂无更多结果，请精确筛选条件'
        })
        if (lockEvent) {
          this.userEventTimeList = {}
        }
        this.refreshDateListShow()
        return new Promise(resolve => {
          resolve(0)
        })
      }

      if (Object.keys(this.userEventTimeList).length == 0) {
        this.userEventTimeList = res.data.list
        this.$message({
          offset: 60,
          type: 'success',
          message: res.msg
        })
        this.refreshDateListShow()
        return new Promise(resolve => {
          resolve(1)
        })
      }
      if (lockEvent) {
        this.userEventTimeList = res.data.list
      } else {
        for (const v of Object.keys(res.data.list)) {
          if (this.userEventTimeList[v]) {
            this.userEventTimeList[v].push(...res.data.list[v])
          } else {
            this.userEventTimeList[v] = res.data.list[v]
          }
        }
      }

      this.$message({
        offset: 60,
        type: 'success',
        message: res.msg
      })
      this.refreshDateListShow()
      return new Promise(resolve => {
        resolve(1)
      })
    },

    async search() {
      const form = JSON.parse(JSON.stringify(this.form))
      form['appid'] = this.$store.state.baseData.EsConnectID

      const res = await UserEventCountList(form)
      if (res.code != 0) {
        this.$message({
          offset: 60,
          type: 'error',
          message: res.msg
        })
        return
      }
      if (res.data.EventPieList == null) res.data.EventPieList = []
      if (res.data.EventLineList == null) res.data.EventLineList = []
      this.chartData.pieData = res.data.EventPieList
      this.chartData.barData = res.data.EventLineList
      this.refreshPage()
    },
    changeLockEvent() {
      this.page = 1
      this.getDescList(true)
    },
    async changeOrderBy(orderBy) {
      this.page = 1
      this.infoForm.sortRule = orderBy
      const status = await this.getDescList()

      if (status == 1) {
        this.page = this.page + 1
      }
    },
    async go() {
      this.page = 1
      this.search()
      const status = await this.getDescList()

      if (status == 1) {
        this.page = this.page + 1
      }
    },
    async filterDateCall(date) {
      this.form.date = date
      this.page = 1
      this.search()
      const status = await this.getDescList()

      if (status == 1) {
        this.page = this.page + 1
      }
    },
    async init() {
      this.form.date = [
        moment().startOf('day').subtract(6, 'days').format('YYYY-MM-DD'),
        moment().startOf('day').subtract(0, 'days').format('YYYY-MM-DD')
      ]

      this.eventOptions = []
      this.form.eventNames = []
      this.eventAttrMap = {}
      const res = await GetConfigs({ 'appid': this.$store.state.baseData.EsConnectID })
      if (res.code != 0) {
        this.$message({
          type: 'error',
          message: res.msg
        })
        return
      }

      this.eventAttrMap = res.data.sys_col

      const event_name_list = res.data.event_name_list

      for (const v of res.data.attributeMap['2']) {
        this.eventAttrMap[v.attribute_name] = v.show_name == '' ? v.attribute_name : v.show_name
      }

      for (const i in event_name_list) {
        const obj = {}
        obj['label'] = event_name_list[i].show_name == '' ? event_name_list[i].event_name : event_name_list[i].show_name
        obj['value'] = event_name_list[i].event_name
        this.eventOptions.push(obj)
        this.form.eventNames.push(event_name_list[i].event_name)
      }
      this.go()
    }
  }
}
</script>

<style scoped src="@/styles/funnel-res.css"/>
<style scoped>
  .el-backtop {
    position: fixed;
    background-color: #fff;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    color: #409eff;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    box-shadow: 0 0 6px rgb(0 0 0 / 12%);
    cursor: pointer;
    z-index: 9999;
  }
</style>
