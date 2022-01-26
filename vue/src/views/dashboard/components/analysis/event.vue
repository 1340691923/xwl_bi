<template>
  <div style="width: 100%;height: 100%;">

    <el-card shadow="hover" style="width: 100%">
      <div slot="header" style="display: flex; align-items: center; justify-content: space-between;">
        <span class="echartBox_title" @click="toRtPage">{{ name }}(事件)</span>
        <div>
          <a-tooltip placement="top" style="cursor: pointer">
            <template slot="title">
              <span>拖移报表</span>
            </template>
            <el-button type="warning" class="drageTag" icon="el-icon-rank" circle />
          </a-tooltip>
          <!--<a-tooltip placement="top" style="cursor: pointer">
            <template slot="title">
              <span>切换成表格或者图表</span>
            </template>
            <el-button type="primary" :icon="totable?'el-icon-s-data':'el-icon-s-grid'" circle @click="changeViewType" />
          </a-tooltip>-->
          <a-tooltip placement="top" style="cursor: pointer">
            <template slot="title">
              <span>刷新</span>
            </template>
            <el-button icon="el-icon-refresh" circle @click="go" />
          </a-tooltip>
          <a-tooltip placement="top" style="cursor: pointer">
            <template slot="title">
              <span>下载表数据</span>
            </template>
            <el-button type="success" icon="el-icon-download" circle @click="download" />
          </a-tooltip>
        </div>
      </div>
      <a-spin tip="计算中..." :spinning="spinning">
        <div class="spin-content">
          <event-result
            v-if="eventResShow"
            :ref="getRef"
            v-model="form.date"
            empty-text="暂无结果，请调整查询条件"
            :window-time-format="form.windowTimeFormat"
            :event-res="eventRes"
            @changeWindowTime="changeWindowTime"
            @go="go"
          />
        </div>
      </a-spin>
    </el-card>
  </div>
</template>

<script>
import { EventList } from '@/api/analysis'
import moment from 'moment'

export default {
  name: 'Event',
  components: {
    'EventResult': () => import('@/views/behavior-analysis/components/EventResult')
  },
  props: {
    name: {
      type: String,
      default: ''
    },
    data: {
      type: Object,
      default: {}
    },
    id: {
      type: String,
      default: ''
    },
    filterDate: {
      type: Array,
      default: []
    }
  },
  data() {
    return {
      eventRes: [],
      spinning: false,
      form: JSON.parse(this.data),
      eventResShow: true,
      tocharts: true,
      totable: false
    }
  },
  computed: {
    getRef() {
      return 'eventRes' + this.id
    }
  },

  watch: {
    'filterDate': {
      immediate: true,
      handler() {
        if (this.filterDate.length > 0) {
          this.form.date = this.filterDate
          Vue.set(this.form, 'date', this.filterDate)
          this.go()
        }
      }
    }
  },
  beforeMount() {
    if (this.filterDate.length == 0) {
      this.form.date = [
        moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD'),
        moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD')
      ]
    }
    this.init()
  },
  methods: {
    changeWindowTime(windowTime) {
      this.form.windowTimeFormat = windowTime
      this.go()
    },
    download() {
      this.$refs[this.getRef].download(this.name)
    },
    async init() {
      await this.go()
    },
    toRtPage() {
      this.$router.push({ path: '/behavior-analysis/event/' + this.id })
    },
    changeViewType() {
      this.tocharts = !this.tocharts
      this.totable = !this.totable
      this.eventResShow = false
      this.$nextTick(() => {
        this.eventResShow = true
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

<style scoped>

.spin-content {
  min-height: 500px;
}
</style>
