<template>
  <div>
    <el-card class="box-card">
      <div style="height: 50px;line-height: 50px;display: flex;align-items: center;justify-content: left">
        <a-tooltip placement="right" style="cursor: pointer">
          <template slot="title">
            <span>添加distinctId为白名单（仅适用于客户端SDK），并开启Debug模式。设备ID可通过查看客户端日志获取。该模式的作用及使用过程请参考官网文档对应SDK页面中的Debug模式说明部分。</span>
          </template>
          <span class="title_xwl" style="color: #202d3f">Debug模式&nbsp<a-icon type="question-circle"/>
          </span>
        </a-tooltip>

      </div>
      <div
        style="height: 60px;line-height: 50px;display: flex;align-items: center;justify-content: space-between;border-bottom: 1px solid #f0f2f5">
        <div>
          <a-button @click="openDeviceDialog">查看distinctId列表</a-button>
          <template v-if="linkedID!=''">
            <el-tag>{{ linkedID }}</el-tag>
            <el-button type="warning" size="mini" @click="closeWs">关闭连接</el-button>
          </template>

        </div>
        <div>
          <a-input-search v-model="input" allow-clear placeholder="请输入搜索（如distinctId,事件名...）" style="width: 200px"
                          @change="search"/>
        </div>
      </div>
      <el-table border v-if="showList"
                v-loading="connectLoading"
                :data="filterList"
                stripe
                style="width: 100%"
      >
        <el-table-column
          prop="report_time"
          label="上报时间"
          align="center"
          width="150"
        />

        <el-table-column
          prop="data_name"
          label="	数据名称"
          align="center"
          width="100"
        />

        <el-table-column label="数据判断" width="200" align="center">
          <template slot-scope="scope">
            <el-tag :type="scope.row.error_reason == undefined?'success':'danger'">
              {{ scope.row.data_judge }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column
          prop="error_reason"
          label="错误原因"
          align="center"
          width="200"
        />

        <el-table-column label="数据明细" align="center">
          <template slot-scope="scope">
            <div v-if="!scope.row.isFormatData">
              {{ scope.row.report_data }}
              <a-tooltip placement="right" style="cursor: pointer">
                <template slot="title">
                  <span>格式化数据</span>
                </template>
                <a-button type="link" icon="eye" @click="lookData(scope.row.index,true)"/>
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
                <a-button style="color: red" type="link" icon="eye-invisible" @click="lookData(scope.row.index,false)"/>
              </a-tooltip>
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
      <el-dialog v-if="dialogVisible" width="80%" :visible.sync="dialogVisible" title="测试distinctId管理"
                 @close="dialogVisible = false">
        <el-card class="app-container">
          <div class="filter-container">
            <el-tag class="filter-item">distinctId:</el-tag>
            <el-input v-model="deviceId" style="width: 300px" placeholder="请输入distinctId" class="filter-item"/>
            <el-tag class="filter-item">备注:</el-tag>
            <el-input v-model="remark" style="width: 200px" placeholder="请输入备注" class="filter-item"/>
            <el-button class="filter-item" type="primary" @click="addDeviceId">添加</el-button>
          </div>
          <el-table border v-if="showDebugDeviceList"
                    :data="debugDeviceList"
                    stripe
                    style="width: 100%"
          >
            <el-table-column
              label="测试distinctId"
              align="center"
              width="300"
            >
              <template slot-scope="scope">
                {{ scope.row.device_id }}
              </template>
            </el-table-column>
            <el-table-column
              label="备注"
              align="center"
              width="200"
            >
              <template slot-scope="scope">
                {{ scope.row.remark }}
              </template>
            </el-table-column>
            <el-table-column
              label="创建时间"
              align="center"
            >
              <template slot-scope="scope">
                {{ scope.row.create_time }}
              </template>
            </el-table-column>
            <el-table-column
              fixed="right"
              label="操作"
              width="450"
              align="center"
            >
              <template slot-scope="scope">
                <el-button size="mini" type="success" icon="el-icon-link" v-if="scope.row.device_id!=linkedID"
                           @click="initWs(scope.row.device_id)">连接设备
                </el-button>
                <el-button size="mini" type="warning" icon="el-icon-link" v-else @click="closeWs()">断开连接</el-button>
                <el-button size="mini" type="danger" icon="el-icon-close" @click="delteDeviceId(scope.row.device_id)">
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-dialog>
    </el-card>

  </div>
</template>

<script>

import {AddDebugDeviceID, DebugDeviceIDList, DelDebugDeviceID} from '@/api/realdata'
import {filterData} from '@/utils/table'
import {getToken} from '@/utils/auth'

export default {

  name: 'Track',
  components: {
    JsonEditor: () => import('@/components/JsonEditor/index')
  },
  data() {
    return {
      showList: true,
      showDebugDeviceList: true,
      dialogVisible: false,
      input: '',
      connectLoading: false,
      time: null,
      total: 0,
      limit: 10,
      page: 1,
      pageshow: true,
      list: [],
      trueList: [],
      debugDeviceList: [],
      deviceId: '',
      remark: '',
      linkedID: '',
      socket: null,
    }
  },
  computed: {
    filterList() {
      var table = this.list.slice((this.page - 1) * this.limit, this.page * this.limit)
      this.refreshPage()
      return table
    },
    path() {
      let api = process.env.VUE_APP_BASE_API

      if(api == ""){

        return "ws://"+ window.location.host+"/ws"

      }

      return api.replace("http", "ws") + "/ws"
    }
  },
  mounted() {
    this.search()
  },
  destroyed() {
    this.closeWs()
  },
  methods: {
    lookData(index, typ) {
      for (const i in this.list) {
        if (this.list[i].index == index) {
          this.list[i].isFormatData = typ
        }
      }
    },
    async delteDeviceId(deviceId) {
      const res = await DelDebugDeviceID({'appid': this.$store.state.baseData.EsConnectID, deviceID: deviceId})
      if (res.code != 0) {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'error',
          message: res.msg
        })
        return
      } else {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'success',
          message: res.msg
        })
        await this.deviceList()
      }
      if (linkedID != '' && deviceId == linkedID) {
        this.closeWs()
      }
    },
    async addDeviceId() {
      const res = await AddDebugDeviceID({
        'appid': this.$store.state.baseData.EsConnectID,
        deviceID: this.deviceId,
        remark: this.remark
      })
      if (res.code != 0) {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'error',
          message: res.msg
        })
        return
      } else {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'success',
          message: res.msg
        })
        await this.deviceList()
      }
    },
    async openDeviceDialog() {
      await this.deviceList()
      this.deviceId = ''
      this.remark = ''
      this.dialogVisible = true
    },
    async deviceList() {
      this.connectLoading = true
      const res = await DebugDeviceIDList({'appid': this.$store.state.baseData.EsConnectID})
      this.connectLoading = false
      if (res.code != 0) {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'error',
          message: res.msg
        })
        return
      }

      let list = res.data.list
      this.debugDeviceList = []
      if (list != null) {
        this.debugDeviceList = list
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
    closeWs() {
      if (this.socket != null) {
        this.socket.close()
        this.socket = null
      }
      this.linkedID = ''
      this.list = []
      this.showDebugDeviceList = false
      this.showList = false
      this.$nextTick(() => {
        this.showDebugDeviceList = true
        this.showList = true
      })
    },
    initWs(uuid) {
      this.closeWs()
      if (this.linkedID != "" && uuid != this.linkedID) {
        this.list = []
      }

      if (typeof (WebSocket) === "undefined") {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'error',
          message: "您的浏览器不支持socket"
        })
      } else {
        // 实例化socket
        this.socket = new WebSocket(this.path)
        // 监听socket连接
        this.socket.onopen = () => {
          this.send(JSON.stringify({uuid: uuid, token: getToken()}))
          this.$message({
            showClose: true,
            offset: 60,
            type: 'success',
            message: "socket连接成功"
          })
        }
        // 监听socket错误信息
        this.socket.onerror = (err) => {
          console.log("onerror", err)
          this.$message({
            showClose: true,
            offset: 60,
            type: 'error',
            message: "socket连接失败",
          })
        }
        // 监听socket消息
        this.socket.onmessage = (msg) => {
          let res = JSON.parse(msg.data)
          if (res.code == 1) {
            let data = JSON.parse(res.data)
            data['isFormatData'] = false
            data['index'] = this.list.length
            data['dataFormat'] = JSON.stringify(JSON.parse(data['report_data']), null, '\t')
            this.list.unshift(data)
            this.search()
          }
        }
        this.socket.onclose = (res) => {
          if (res.code == 1006) {
            setTimeout(() => {
              console.log("进行重连")
              this.initWs(uuid)
            }, 2000)
          }
          console.log("socket已经关闭")
        }
        this.linkedID = uuid
      }
    },
    send(data) {
      this.socket.send(data)
    },
    async search() {
      let list = filterData(this.list, this.input.trim())
      this.list = list
    }
  }
}
</script>

<style scoped src="@/styles/debug.css"/>
