<template>
  <div>
    <el-card class="box-card">
      <div style="height: 50px;line-height: 50px;display: flex;align-items: center;justify-content: space-between">
        <a-tooltip placement="right" style="cursor: pointer">
          <template slot="title">
            <span>实时入库：展示最近1000条入库数据；错误数据：按错误原因展示最近1000条统计数据，按10分钟/次刷新统计结果。</span>
          </template>
          <span class="title_xwl" style="color: #202d3f">
            实时数据
            <a-icon type="question-circle" />
          </span>
        </a-tooltip>
        <template>
          <el-button
            v-if="refreshtActiveName == 'realTime2Es'"
            v-loading="downloading"
            size="mini"
            type="text"
            @click="downloadData"
          >下载数据 <i class="el-icon-download el-icon--right" /></el-button>
        </template>
      </div>
      <div
        style="height: 50px;line-height: 50px;display: flex;align-items: center;justify-content: space-between;border-bottom: 1px solid #f0f2f5"
      >
        <div>
          <el-tabs v-model="refreshtActiveName" size="small" @tab-click="handleClick">
            <el-tab-pane label="实时入库" name="realTime2Es" />
            <el-tab-pane label="错误数据" name="failData" />
          </el-tabs>
        </div>
        <div>
          <template>
            <div>
              <date :dates="date" style="float: left" @changeDate="changeDate" />
              <div style="float: left;margin-left: 10px">
                轮循/S：
                <a-input-number
                  v-model="timeSecend"
                  type="number"
                  placeholder="轮循"
                  style="width: 60px"
                  @change="changeRealDataShow"
                />
                <a-button type="link" :icon="loopStatus?'pause':'play-circle'" @click="changeLoopStatus" />
                <a-input-search v-model="searchKw" placeholder="服务器过滤" allow-clear style="width: 150px" />
                <a-input-search v-model="input" placeholder="客户端过滤" allow-clear style="width: 150px" />
              </div>
            </div>
          </template>
        </div>
      </div>
      <template v-if="refreshtActiveName == 'realTime2Es'">
        <real-data2-es
          v-if="showRealData"
          ref="realTime2Es"
          :date="dateFormart"
          :search-kw="searchKw"
          :input="input"
          :time-secend="timeSecend"
          :pause="pause"
          @downloadFinish="downloadFinish"
        />
      </template>
      <template v-if="refreshtActiveName == 'failData'">
        <fail-data v-if="showRealData" :input="input" :time-secend="timeSecend" :pause="pause" />
      </template>
    </el-card>
  </div>
</template>

<script>

import moment from 'moment'

export default {
  name: 'RealTime',
  components: {
    'RealData2Es': () => import('@/views/manager/components/realData2Es'),
    'FailData': () => import('@/views/manager/components/FailData'),
    Date: () => import('@/components/Date')
  },
  data() {
    return {
      date: [
        moment().startOf('day').format('YYYY-MM-DD HH:mm:ss'),
        moment().endOf('day').format('YYYY-MM-DD HH:mm:ss')
      ],
      downloading: false,
      timeSecend: 30,
      input: '',
      loopStatus: true,
      pause: false,
      activeName: 'realTime2Es',
      showRealData: true,
      searchKw: ''
    }
  },
  computed: {
    dateFormart() {
      if (this.date.length == 0) return ''
      return this.date.join(',')
    },
    refreshtActiveName: {
      get() {
        if (this.$store.state.baseData.activeName == '') {
          return this.activeName
        }
        this.activeName = this.$store.state.baseData.activeName
        return this.$store.state.baseData.activeName
      },
      set(val) {
        this.$store.dispatch('baseData/SETActiveName', val)
        this.activeName = val
      }
    }
  },
  mounted() {

  },
  methods: {
    changeDate(date) {
      this.date = date
    },
    downloadData() {
      this.downloading = true
      this.$refs[this.refreshtActiveName].download2Excel()
    },
    changeRealDataShow() {
      this.showRealData = false
      this.$nextTick(() => {
        this.showRealData = true
      })
    },
    downloadFinish() {
      this.downloading = false
    },
    changeLoopStatus() {
      this.loopStatus = !this.loopStatus
      this.pause = !this.pause
      if (this.pause) {
        this.$message({
          offset: 60,

          type: 'success',
          message: '已暂停轮循'
        })
      } else {
        this.$message({
          offset: 60,

          type: 'success',
          message: '已开始轮循'
        })
        this.changeRealDataShow()
      }
    }
  }
}
</script>

<style scoped src="@/styles/real-time.css"/>
