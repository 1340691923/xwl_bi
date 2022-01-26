<template>
  <div>
    <el-card class="box-card">
      <div style="height: 50px;line-height: 50px;display: flex;align-items: center;justify-content: left">
        <a-tooltip placement="right" style="cursor: pointer">
          <template slot="title">
            <span>查看最近7日项目内事件或用户属性的数据接收情况，对事件进行禁用/启用管理。</span>
          </template>
          <span class="title_xwl" style="color: #202d3f">上报统计 <a-icon
            type="question-circle"
          />
          </span>
        </a-tooltip>

      </div>
      <div
        style="height: 60px;line-height: 50px;display: flex;align-items: center;justify-content: space-between;border-bottom: 1px solid #f0f2f5"
      >
        <div>
          <date :dates="date" @changeDate="changeDate" />
        </div>
        <div>
          <a-input-search v-model="input" placeholder="请输入搜索" style="width: 200px" @change="search" />
        </div>
      </div>
      <el-table
        v-loading="connectLoading"
        border
        :data="filterList"
        stripe
        style="width: 100%"
      >
        <el-table-column
          prop="data_name"
          label="数据名称"
          align="center"
          width="250"
        />

        <el-table-column
          prop="show_name"
          label="显示名"
          align="center"
          width="200"
        />
        <el-table-column
          prop="received_count"
          label="已接收"
          align="center"
          width="250"
        />
        <el-table-column
          prop="succ_count"
          label="已入库"
          align="center"
          width="200"
        />
        <el-table-column
          prop="fail_count"
          label="入库失败"
          align="center"
          width="200"
        />

        <el-table-column
          fixed="right"
          label="操作"
          width="200"
          align="center"
        >
          <template slot-scope="scope">
            <a v-if="scope.row.fail_count != 0" style="color: #6bb8ff" @click="openFailDesc(scope.row)">错误详情</a>
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
      <el-dialog
        v-if="dialogVisible"
        width="80%"
        :visible.sync="dialogVisible"
        :title="failTitleList[0].data_name"
        @close="dialogVisible = false"
      >
        <el-card>
          <el-table
            v-loading="connectLoading"
            border
            :data="failTitleList"
            stripe
            style="width: 100%"
          >
            <el-table-column
              prop="start_time"
              label="开始时间"
              align="center"
              width="250"
            />

            <el-table-column
              prop="end_time"
              label="结束时间"
              align="center"
              width="200"
            />
            <el-table-column
              prop="received_count"
              label="已接收"
              align="center"
              width="250"
            />
            <el-table-column
              prop="succ_count"
              label="已入库"
              align="center"
              width="200"
            />
            <el-table-column
              prop="fail_count"
              label="入库失败"
              align="center"
            />

          </el-table>
        </el-card>
        <el-card style="margin-top: 10px">

          <el-table
            border
            :data="failDescList"
            stripe
            style="width: 100%"
          >
            <el-table-column
              prop="error_reason"
              label="错误原因"
              align="center"
              width="250"
            />

            <el-table-column
              prop="count"
              label="错误条数"
              align="center"
              width="200"
            />

            <el-table-column label="抽样示例" align="center">
              <template slot-scope="scope">
                <div v-if="!scope.row.isFormatData">
                  {{ scope.row.report_data }}
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
                    <a-button
                      style="color: red"
                      type="link"
                      icon="eye-invisible"
                      @click="lookData(scope.row.index,false)"
                    />
                  </a-tooltip>
                </div>
              </template>
            </el-table-column>

          </el-table>
        </el-card>
      </el-dialog>
    </el-card>

  </div>
</template>

<script>

import { EventFailDesc, ReportCount } from '@/api/realdata'
import { filterData } from '@/utils/table'

export default {

  name: 'Track',
  components: {
    Date: () => import('@/components/Date'),
    JsonEditor: () => import('@/components/JsonEditor/index')
  },
  data() {
    return {
      dialogVisible: false,
      input: '',
      connectLoading: false,
      time: null,
      total: 0,
      limit: 100,
      page: 1,
      pageshow: true,
      list: [],
      date: [],
      trueList: [],
      failDescList: [],
      failTitleList: []
    }
  },
  computed: {
    filterList() {
      var table = this.list.slice((this.page - 1) * this.limit, this.page * this.limit)
      this.refreshPage()
      return table
    }
  },
  mounted() {
    this.date.push(
      this.$moment().startOf('day').format('YYYY-MM-DD HH:mm:ss'), this.$moment().endOf('day').format('YYYY-MM-DD HH:mm:ss')
    )
    this.search()
  },
  methods: {
    lookData(index, typ) {
      for (const i in this.list) {
        if (this.failDescList[i].index == index) {
          this.failDescList[i].isFormatData = typ
        }
      }
    },
    async openFailDesc(row) {
      const res = await EventFailDesc({
        'appid': this.$store.state.baseData.EsConnectID,
        'start_time': this.date[0],
        'end_time': this.date[1],
        'data_name': row['data_name']
      })
      if (res.code != 0) {
        this.$message({
          offset: 60,

          type: 'error',
          message: res.msg
        })
        return
      }
      this.failTitleList = []
      row['start_time'] = this.date[0]
      row['end_time'] = this.date[1]
      this.failTitleList.push(row)
      this.failDescList = []

      for (const index in res.data.list) {
        const v = res.data.list[index]
        v['index'] = index
        v['isFormatData'] = false
        v['dataFormat'] = JSON.stringify(JSON.parse(v['report_data']), null, '\t')
        this.failDescList.push(v)
      }

      this.dialogVisible = true
    },
    changeDate(val) {
      this.date = val
      this.search()
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
    async search() {
      if (this.date.length < 2) {
        this.$message({
          offset: 60,

          type: 'error',
          message: '请先选择筛选时间'
        })
        return
      }
      this.connectLoading = true
      const res = await ReportCount({
        'appid': this.$store.state.baseData.EsConnectID,
        'start_time': this.date[0],
        'end_time': this.date[1]
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
        if (res.data.list == null) {
          res.data.list = []
        }
        let list = res.data.list
        list = filterData(list, this.input.trim())
        this.total = list.length
        this.list = list
      }
    }
  }
}
</script>

<style scoped src="@/styles/track-data.css"/>
