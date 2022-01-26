<template>
  <div>

    <el-table
      v-loading="connectLoading"
      border
      max-height="700"
      :data="filterList"
      style="width: 100%"
    >
      <el-table-column
        label="序号"
        align="center"
        fixed
        width="50"
      >
        <template slot-scope="scope">
          {{ scope.$index + 1 }}
        </template>
      </el-table-column>

      <el-table-column label="采样日期" width="180" align="center" prop="year" sortable />
      <el-table-column label="采样时间" width="180" align="center">
        <template slot-scope="scope">
          {{ scope.row.start_minute }}-{{ scope.row.end_minute }}
        </template>
      </el-table-column>
      <el-table-column label="错误条数" width="180" align="center" prop="count" sortable />
      <el-table-column label="错误处理" width="180" align="center" prop="error_handling" />
      <el-table-column label="错误类型" width="180" align="center" prop="report_type" />
      <el-table-column label="错误原因" align="center">
        <template slot-scope="scope">
          <div>
            <a style="color: #6bb8ff" @click="lookData(scope.row)">{{ scope.row.error_reason }}</a>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-if="pageshow"
      class="pagination-container"
      :current-page="page"
      :page-sizes="[5,10, 20, 30, 50,100,150,200,500,1000]"
      :page-size="limit"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total"
      @size-change="handleSizeChange"
      @current-change="handlePageChange"
    />
    <el-drawer
      ref="drawer"
      title="抽样数据"
      :before-close="drawerHandleClose"
      :visible.sync="drawerShow"

      direction="rtl"
      close-on-press-escape
      destroy-on-close
      size="50%"
    >

      <json-editor
        v-if="drawerShow"
        v-model="descData"
        font-size="15"
        height="800"
        styles="width: 100%"
        :read="true"
        title="抽样数据"
      />
    </el-drawer>

  </div>
</template>

<script>
import { FailDataDesc, FailDataList } from '@/api/realdata'
import { filterData } from '@/utils/table'

export default {
  name: 'RealData2Es',
  components: {
    'JsonEditor': () => import('@/components/JsonEditor/index')
  },
  props: {
    input: {
      type: String,
      default: ''
    },
    timeSecend: {
      type: Number,
      default: 5
    },
    pause: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      descData: '{}',
      drawerShow: false,
      connectLoading: false,
      time: null,
      total: 0,
      limit: 10,
      page: 1,
      pageshow: true,
      list: [],
      trueList: []
    }
  },
  computed: {
    filterList() {
      var table = this.list.slice((this.page - 1) * this.limit, this.page * this.limit)
      this.refreshPage()
      return table
    }
  },
  watch: {
    pause(newV, oldV) {
      if (newV == true) {
        this.cleanTimer()
      }
    },
    input(newV, oldV) {
      let list = this.trueList
      list = filterData(list, newV.trim())
      this.total = list.length
      this.list = list
    }
  },
  beforeDestroy() {
    // 清除定时器
    clearInterval(this.time)
    this.time = null
  },
  mounted() {
    this.searchData()
    if (!this.pause) {
      this.startLoop()
    } else {
      clearInterval(this.time)
      this.time = null
    }
  },
  methods: {
    cleanTimer() {
      clearInterval(this.time)
      this.time = null
    },
    async lookData(row) {
      const form = {
        start_time: `${row['year']} ${row['start_minute']}:00`,
        end_time: `${row['year']} ${row['end_minute']}:00`,
        appid: this.$store.state.baseData.EsConnectID,
        error_reason: row['error_reason'],
        error_handling: row['error_handling'],
        report_type: row['report_type']
      }
      const res = await FailDataDesc(form)
      if (res.code != 0) {
        this.$message({
          offset: 60,

          type: 'error',
          message: res.msg
        })
        return
      } else {
        this.descData = JSON.stringify(JSON.parse(res.data.data), null, '\t')
        this.drawerShow = true
      }
    },
    handleSizeChange(v) {
      this.limit = v
      this.refreshPage()
    },
    handlePageChange(v) {
      this.page = v
      this.refreshPage()
    },
    refreshPage() {
      this.pageshow = false
      this.total = this.list.length
      this.$nextTick(() => {
        this.pageshow = true
      })
    },
    async searchData() {
      this.connectLoading = true
      const res = await FailDataList({ 'appid': this.$store.state.baseData.EsConnectID })
      this.connectLoading = false
      if (res.code != 0) {
        this.$message({
          offset: 60,

          type: 'error',
          message: res.msg
        })
        return
      } else {
        let list = []
        let index = 0
        for (const v of res.data.list) {
          v['index'] = index
          list.push(v)
          index++
        }
        list = filterData(list, this.input.trim())

        this.total = list.length
        this.list = list
        this.trueList = list
      }
    },

    startLoop() {
      this.time = setInterval(() => {
        this.searchData()
      }, this.timeSecend * 1000)
    },
    drawerHandleClose() {
      this.drawerShow = false
    }
  }
}
</script>

<style scoped>

</style>
