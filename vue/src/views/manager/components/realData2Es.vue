<template>
  <div>

    <el-table
      v-loading="connectLoading"
      border
      :data="filterList"
      style="width: 100%"
      max-height="700"
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
      <el-table-column label="上报时间" width="180" align="center" prop="create_time" sortable />
      <el-table-column label="数据名称" width="180" align="center" prop="event_name" />
      <el-table-column label="数据明细" align="center">
        <template slot-scope="scope">
          <div v-if="!scope.row.isFormatData">
            <span v-html="scope.row.data " />
            <a-tooltip placement="right" style="cursor: pointer">
              <template slot="title">
                <span>格式化数据</span>
              </template>
              <a-button type="link" icon="eye" @click="lookData(scope.row.index,true)" />
            </a-tooltip>
          </div>
          <div v-else>
            <json-editor
              v-model="scope.row.dataFormat"
              font-size="15"
              height="400"
              class="req-body"
              styles="width: 100%"
              :read="true"
              title="上报数据"
            />
            <a-tooltip placement="right" style="cursor: pointer">
              <template slot="title">
                <span>还原数据</span>
              </template>
              <a-button style="color: red" type="link" icon="eye-invisible" @click="lookData(scope.row.index,false)" />
            </a-tooltip>
          </div>
        </template>
      </el-table-column>

    </el-table>
    <el-pagination
      v-if="pageshow"
      class="pagination-container"
      :current-page="page"
      :page-sizes="[10, 20, 30, 50,100,150,200,500,1000]"
      :page-size="limit"
      layout="total, sizes, prev, pager, next, jumper"
      :total="total"
      @size-change="handleSizeChange"
      @current-change="handlePageChange"
    />
  </div>
</template>

<script>
import { List } from '@/api/realdata'
import { filterData } from '@/utils/table'
import { Message } from 'element-ui'

export default {
  name: 'RealData2Es',
  components: {
    'JsonEditor': () => import('@/components/JsonEditor/index')
  },

  props: {
    date: {
      type: String,
      default: ''
    },
    searchKw: {
      type: String,
      default: ''
    },
    input: {
      type: String,
      default: ''
    },
    timeSecend: {
      type: Number,
      default: 30
    },
    pause: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      connectLoading: false,
      time: null,
      total: 0,
      limit: 100,
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
    },
    searchKw(newV, oldV) {
      this.searchData()
    },
    date(newV, oldV) {
      this.searchData()
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
    formatJson(filterVal, jsonData) {
      return jsonData.map(v => filterVal.map(j => v[j]))
    },
    download2Excel() {
      import('@/vendor/Export2Excel').then(excel => {
        const $message = Message.success({
          message: '  数据下载中...',
          type: 'success',
          duration: 0,
          iconClass: 'el-icon-loading'
        })

        const tHeader = [
          '上报时间',
          '数据名称',
          '数据明细'
        ]

        // 需要导出的列名
        const filterVal = [
          'create_time',
          'event_name',
          'data'
        ]

        // 接口获得的所有数据
        const list = this.list
        const data = this.formatJson(filterVal, list)
        excel.export_json_to_excel({
          header: tHeader,
          data,
          filename: '实时入库数据',
          autoWidth: true,
          bookType: 'xlsx'
        })
        this.$emit('downloadFinish')
        $message.close()
      })
    },
    cleanTimer() {
      clearInterval(this.time)
      this.time = null
    },
    lookData(index, typ) {
      for (const i in this.list) {
        if (this.list[i].index == index) {
          this.list[i].isFormatData = typ
        }
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
      const res = await List({
        'appid': this.$store.state.baseData.EsConnectID,
        'searchKw': this.searchKw,
        date: this.date
      })
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
          const _source = {}

          _source['dataFormat'] = JSON.stringify(JSON.parse(v["report_data"]), null, '\t')
          _source['event_name'] = v["event_name"]
          _source['create_time'] = v["create_time"]
          _source['data'] = v["report_data"]
          _source['isFormatData'] = false
          _source['index'] = index
          list.push(_source)
          index++
        }
        list = filterData(list, this.input.trim())
        console.log("list",list)
        this.total = list.length
        this.list = list
        this.trueList = list
      }
    },

    startLoop() {
      this.time = setInterval(() => {
        this.searchData()
      }, this.timeSecend * 1000)
    }

  }
}
</script>

<style scoped>

</style>
