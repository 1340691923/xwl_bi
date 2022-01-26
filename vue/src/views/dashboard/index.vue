<template>
  <div style="display:flex;justify-content:space-between">

    <div class="sider_xwl">
      <div class="top_xwl">
        <manager-dir
          v-model="panelStr"
          style="padding-left: 5px;padding-right: 5px"
          @change="GetPannelList"
        />
        <a-button size="small" type="danger" icon="delete" @click.native="deletePannelVisible = true" />
        <a-dropdown style="margin-left: 10px" :trigger="['click']" placement="bottomRight">
          <a-button size="small" type="primary" icon="plus" />
          <a-menu slot="overlay">
            <a-menu-item key="1" @click="openAddPannel">
              <a-icon type="dashboard" />
              新建看板
            </a-menu-item>
            <a-menu-item key="2" @click="newDir">
              <a-icon type="folder" />
              新建文件夹
            </a-menu-item>
            <a-menu-divider />

          </a-menu>
        </a-dropdown>
      </div>
      <div class="body_xwl">

        <div>
          <a-menu
            v-model="LastSelectKey"
            :open-keys="openKeys"
            mode="inline"
            @openChange="onOpenChange"
          >
            <a-sub-menu key="dashboard">
              <span slot="title">
                <a-icon type="dashboard" />
                <span>我的看板</span>
              </span>

              <a-sub-menu v-for="(v,k,index) in pannelList" :key="v.folder_id">

                <span slot="title">
                  <a-icon type="folder" /><span>{{ v.folder_name }}</span>
                </span>

                <a-menu-item
                  v-for="(v2,k2,index2) in v.childList"
                  :key="v.folder_id.toString()+'.'+v2.pannel_id.toString()"
                  @click="lookPannel(v2.pannel_id,v2.pannel_name,v2.report_tables,v.folder_type == 0,v2.managers,v2.create_by)"
                >
                  <a-dropdown :trigger="['contextmenu']">
                    <div
                      class="dashbordName_xwl"
                      :title="v2.pannel_name"
                    >
                      {{ v2.pannel_name }}
                    </div>
                    <a-menu v-if="v.folder_type == 0" slot="overlay">
                      <a-menu-item key="1" @click="rename(v2.pannel_id)">
                        重命名
                      </a-menu-item>
                      <a-menu-item key="2" @click="openMovePannelDialog(v2.pannel_id,v.folder_id)">
                        移动至
                      </a-menu-item>
                      <a-menu-item key="3" @click="openCopyPannel(v2.report_tables,v2.pannel_name,v.folder_id)">
                        复制看板
                      </a-menu-item>
                      <a-menu-item key="4" @click="deletePannel(v2.pannel_id,v2.pannel_name)">
                        删除看板
                      </a-menu-item>
                    </a-menu>
                  </a-dropdown>
                </a-menu-item>
              </a-sub-menu>
            </a-sub-menu>
          </a-menu>
        </div>
      </div>
    </div>
    <div class="content_xwl">
      <div class="header_xwl" style="background: white">
        <div class="root_xwl">
          <div class="main_xwl">
            <span
              class="title_xwl"
              :title="selectPannel.pannelName"
              style="color: #202d3f"
            >&nbsp;&nbsp;{{ selectPannel.pannelName }}</span>
          </div>
          <div class="actions_xwl">

            <a-button-group>
              <span style="cursor: pointer;border: 1px solid #e4f2ff;    border-radius: 2px;">
                <a-range-picker format="YYYY-MM-DD" :ranges="dataRange" @change="filterDateCall">
                  <div style="cursor: pointer">
                    <a-button type="link" class="actions_xwl_btn" icon="clock-circle" />{{ filterDate.join("~") }}
                  </div>
                </a-range-picker>

                <a-button
                  v-if="filterDate.length > 0"
                  type="link"
                  style="color: red"
                  class="actions_xwl_btn"
                  icon="close-circle"
                  @click="cleanDateFiter"
                /></span>
              <span style="border-radius: 2px;">
                <a-button type="link" class="actions_xwl_btn" icon="reload" @click="loadReportTables" />
                <a-button
                  v-if="selectPannel.canShare"
                  type="link"
                  class="actions_xwl_btn"
                  icon="share-alt"
                  @click="openManagerForm"
                />
                <a-button
                  v-if="selectPannel.canShare"
                  type="link"
                  class="actions_xwl_btn"
                  icon="drag"
                  @click="openAddPannelRt"
                />
              </span>
            </a-button-group>
          </div>
        </div>
      </div>

      <div
        v-if="selectPannel.pannelId == 0"
        style="background: white !important;padding: 40px;width: 300px;height: 300px; text-align: center;margin: 0px auto;position: relative;top: 30%"
      >
        <a-empty>
          <span slot="description">请点击左侧面板查看报表</span>
        </a-empty>
      </div>
      <template v-else-if="selectPannel.reportTables.length == 0">
        <div
          v-if="selectPannel.canShare "
          style="background: white !important;padding: 40px;width: 300px;height: 300px; text-align: center;margin: 0px auto;position: relative;top: 30%"
        >
          <a-empty>
            <span slot="description">尚未添加报表</span>
            <a-button type="primary" @click="openAddPannelRt">
              添加报表
            </a-button>
          </a-empty>
        </div>
        <div
          v-else
          style="background: white !important;padding: 40px;width: 300px;height: 300px; text-align: center;margin: 0px auto;position: relative;top: 30%"
        >
          <a-empty>
            <span slot="description">该用户未添加任何报表</span>
          </a-empty>
        </div>
      </template>

      <div v-else style="position: relative; height: calc(100% - 56px);overflow-y: auto;">
        <template v-if="echartBoxShow">
          <draggable
            v-model="selectPannel.reportTables"
            :handle="selectPannel.canShare?'.drageTag':'cant_drag'"
            style="width: 100%;flex-wrap:wrap;display:flex; flex-direction: row;padding-top: 20px;padding-left: 20px"
            animation="1000"
            @end="onEnd"
          >
            <div v-for="(v,k,index) in selectPannel.reportTables" class="echartBox" :style="echartsBoxStyle(v)">
              <a-dropdown :trigger="['contextmenu']">
                <event
                  v-if="rtConfig[v].rt_type == 1"
                  :id="v"
                  :key="index"
                  :filter-date="filterDate"
                  :data="rtConfig[v].data"
                  :name="rtConfig[v].name"
                />
                <funnel
                  v-if="rtConfig[v].rt_type == 3"
                  :id="v"
                  :key="index"
                  :filter-date="filterDate"
                  :data="rtConfig[v].data"
                  :name="rtConfig[v].name"
                />
                <retention
                  v-if="rtConfig[v].rt_type == 2"
                  :id="v"
                  :key="index"
                  :filter-date="filterDate"
                  :data="rtConfig[v].data"
                  :name="rtConfig[v].name"
                />
                <trace
                  v-if="rtConfig[v].rt_type == 4"
                  :id="v"
                  :key="index"
                  :filter-date="filterDate"
                  :data="rtConfig[v].data"
                  :name="rtConfig[v].name"
                />
                <user-attr
                  v-if="rtConfig[v].rt_type == 5"
                  :id="v"
                  :key="index"
                  :data="rtConfig[v].data"
                  :name="rtConfig[v].name"
                />
                <a-menu slot="overlay">
                  <a-menu-item key="1" @click="changeBigScreen(v)">
                    {{ reportTablesBigIdMap.hasOwnProperty(v) ? '变成小屏' : '变成大屏' }}
                  </a-menu-item>
                </a-menu>
              </a-dropdown>
            </div>
          </draggable>
        </template>
      </div>
    </div>
    <el-dialog :close-on-click-modal="false" :visible.sync="addPannelVisible" title="添加看板" width="80%" append-to-body>
      <el-form :model="addPannelForm" label-width="120px" label-position="left">
        <el-form-item label="看板名称:">
          <el-input v-model="addPannelForm.pannel_name" style="width:500px" />
        </el-form-item>
        <el-form-item label="添加至:">
          <el-select v-model="addPannelForm.folder_id" placeholder="请选择文件夹" filterable>

            <el-option
              v-for="(v,k,index) in pannelList"
              v-if="v.folder_type == 0"
              :key="index"
              :label="v.folder_name"
              :value="v.folder_id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div style="text-align:right;">
        <el-button type="danger" icon="el-icon-close" @click="addPannelVisible = false">返回</el-button>
        <el-button type="primary" icon="el-icon-plus" @click="addPannel">添加</el-button>
      </div>
    </el-dialog>

    <el-dialog
      :close-on-click-modal="false"
      :visible.sync="copyPannelVisible"
      title="复制看板（复制看板将同时复制看板内的报表）"
      width="80%"
      append-to-body
    >
      <el-form :model="copyPannelForm" label-width="120px" label-position="left">
        <el-form-item label="看板名称:">
          <el-input v-model="copyPannelForm.pannel_name" style="width:500px" />
        </el-form-item>
        <el-form-item label="复制至:">
          <el-select v-model="copyPannelForm.folder_id" placeholder="请选择文件夹" filterable>
            <el-option
              label="请先建立自己的文件夹"
              :value="Number(-1)"
            />
            <el-option
              v-for="(v,k,index) in pannelList"
              v-if="v.folder_type == 0"
              :key="index"
              :label="v.folder_name"
              :value="v.folder_id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div style="text-align:right;">
        <el-button type="danger" icon="el-icon-close" @click="copyPannelVisible = false">返回</el-button>
        <el-button type="primary" icon="el-icon-plus" @click="copyPannel">复制至</el-button>
      </div>
    </el-dialog>

    <el-dialog :close-on-click-modal="false" :visible.sync="movePannelVisible" title="移动至" width="80%" append-to-body>
      <el-form :model="movePannelForm" label-width="120px" label-position="left">
        <el-form-item label="移动至:">
          <el-select v-model="movePannelForm.folder_id" placeholder="请选择文件夹" filterable>

            <el-option
              v-for="(v,k,index) in pannelList"
              v-if="v.folder_type == 0"
              :key="index"
              :label="v.folder_name"
              :value="v.folder_id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div style="text-align:right;">
        <el-button type="danger" icon="el-icon-close" @click="movePannelVisible = false">返回</el-button>
        <el-button type="primary" icon="el-icon-check" @click="movePannel">移动</el-button>
      </div>
    </el-dialog>

    <el-dialog
      :close-on-click-modal="false"
      :visible.sync="deletePannelVisible"
      title="删除文件夹(文件夹的面板也一同会删除)"
      width="50%"
      append-to-body
    >
      <a-tree
        v-model="checkedKeys"
        checkable
        :tree-data="treeData"
      />
      <div style="text-align:right;">
        <el-button type="danger" icon="el-icon-close" @click="deletePannelVisible = false">返回</el-button>
        <el-button type="primary" icon="el-icon-check" @click="deleteDir">删除</el-button>
      </div>
    </el-dialog>

    <el-dialog
      width="70%"
      :close-on-click-modal="false"
      :visible.sync="managerFormdialogVisible"
      title="分享面板给其他成员"
      @close="managerFormdialogVisible = false"
    >
      <el-form :model="selectPannel" label-width="120px" label-position="left">
        <el-form-item label="面板名">
          <el-input v-model="selectPannel.pannelName" readonly placeholder="应用名" />
        </el-form-item>
        <el-form-item label="成员">
          <el-transfer
            v-if="managerFormdialogVisible"
            v-model="selectPannel.managers"
            :titles="['当前应用全部成员', '当前面板成员']"
            :button-texts="['移除成员', '添加成员']"
            filterable
            :filter-method="filterMethod"
            filter-placeholder="请操作成员"
            :data="allUserConfig"
          />

        </el-form-item>
      </el-form>
      <div style="text-align:right;">
        <el-button type="danger" icon="el-icon-close" @click="managerFormdialogVisible = false">返回</el-button>
        <el-button type="primary" icon="el-icon-check" @click="addManager">添加</el-button>
      </div>
    </el-dialog>

    <manager-rt
      v-if="addPannelRtVisible"
      v-model="selectPannel.reportTables"
      :visible="addPannelRtVisible"
      @saveRt="saveRt"
      @close="addPannelRtVisible = false"
    />
    <back-to-top />
  </div>
</template>

<script>
import moment from 'moment'
import {
  CopyPannel,
  DeleteDir,
  DeletePannel,
  GetPannelList,
  MovePannel2Dir,
  NewDir,
  NewPannel,
  Rename,
  RtListByAppid,
  UpdatePannelManager,
  UpdatePannelRt
} from '@/api/pannel'
import { filterData } from '@/utils/table'
import { userList } from '@/api/user'
import draggable from 'vuedraggable'

const dataRange = {
  '今天': [moment().startOf('day'), moment()],
  '昨天': [
    moment().startOf('day').subtract(1, 'days'),
    moment().startOf('day').subtract(1, 'days')
  ],
  '最近一周': [moment().startOf('day').subtract(1, 'weeks'), moment()],
  '最近两周': [moment().startOf('day').subtract(2, 'weeks'), moment()],
  '最近1个月': [moment().startOf('day').subtract(1, 'months'), moment()],
  '最近3个月': [moment().startOf('day').subtract(3, 'months'), moment()],
  '最近半年': [moment().startOf('day').subtract(6, 'months'), moment()],
  '最近1年': [moment().startOf('day').subtract(1, 'years'), moment()]
}

export default {
  name: 'Index',
  components: {
    ManagerDir: () => import('@/views/dashboard/components/ManagerDir'),
    ManagerRt: () => import('@/views/dashboard/components/ManagerRt'),
    Funnel: () => import('@/views/dashboard/components/analysis/funnel'),
    Retention: () => import('@/views/dashboard/components/analysis/retention'),
    Event: () => import('@/views/dashboard/components/analysis/event'),
    UserAttr: () => import('@/views/dashboard/components/analysis/user_attr'),
    Trace: () => import('@/views/dashboard/components/analysis/trace'),
    BackToTop: () => import('@/components/BackToTop/index'),
    draggable
  },
  data() {
    return {
      openKeys: [],
      echartBoxShow: true,
      userConfig: [],
      allUserConfig: [],
      managerFormdialogVisible: false,
      addPannelRtVisible: false,
      selectKey: [],
      selectPannel: {
        create_by: 0,
        pannelId: 0,
        pannelName: '请先选中面板',
        reportTables: [],

        managers: '',
        canShare: false
      },
      reportTablesBigIdMap: {},
      rtConfig: {},
      treeData: [],
      checkedKeys: [],
      deletePannelVisible: false,
      copyPannelForm: {
        folder_id: '',
        pannel_name: '',
        appid: this.$store.state.baseData.EsConnectID,
        report_tables: ''
      },
      copyPannelVisible: false,
      movePannelVisible: false,
      movePannelForm: {
        folder_id: '',
        id: 0
      },
      addPannelForm: {
        folder_id: '',
        pannel_name: ''
      },
      addPannelVisible: false,
      menuTyp: 1,
      pannelList: [],
      filterDate: [],
      panelStr: '',
      dataRange
    }
  },
  computed: {
    LastSelectKey: {
      get() {
        if (this.$store.state.baseData.LastSelectKey == '') {
          return this.selectKey
        }
        this.selectKey = this.$store.state.baseData.LastSelectKey
        return this.$store.state.baseData.LastSelectKey
      },
      set(val) {
        this.$store.dispatch('baseData/SETLastSelectKey', val)
        this.selectKey = val
      }
    }
  },
  watch: {
    'selectPannel.pannelId'(newV, oldV) {
      this.loadReportTables()
    }
  },
  mounted() {
    this.init()
  },
  methods: {
    changeBigScreen(id) {
      if (this.reportTablesBigIdMap.hasOwnProperty(id)) {
        Vue.delete(this.reportTablesBigIdMap, id)
      } else {
        Vue.set(this.reportTablesBigIdMap, id, 1)
      }
    },
    echartsBoxStyle(id) {
      // 如果没有特殊处理 都是宽度默认48%
      if (!this.reportTablesBigIdMap.hasOwnProperty(id)) {
        return {
          width: '50%'
        }
      } else {
        return {
          width: '100%'
        }
      }
    },
    onEnd() {
      this.saveRt(this.selectPannel.reportTables)
    },

    async addManager() {
      const managers = this.selectPannel.managers.join(',')
      const id = this.selectPannel.pannelId
      const appid = this.$store.state.baseData.EsConnectID
      const res = await UpdatePannelManager({ managers, id, appid })
      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })
        return
      } else {
        this.$message({
          type: 'success',
          offset: 60,
          message: res.msg
        })
      }
      this.GetPannelList()
      this.managerFormdialogVisible = false
    },
    async init() {
      await this.rtListByAppid()
      await this.GetPannelList()
    },
    async openManagerForm() {
      this.allUserConfig = []
      const userListRes = await userList({ appid: this.$store.state.baseData.EsConnectID })

      if (userListRes) {
        for (var v of userListRes.data) {
          if (this.selectPannel.create_by != v.id) {
            this.userConfig[v.id] = v.realname
            const obj = {
              label: v.realname,
              key: v.id.toString(),
              disabled: false
            }

            this.allUserConfig.push(
              obj
            )
          }
        }
      }

      this.managerFormdialogVisible = true
    },
    filterMethod(query, item) {
      return item.label.indexOf(query) > -1
    },
    async rtListByAppid() {
      const res = await RtListByAppid({ appid: this.$store.state.baseData.EsConnectID })
      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })
        return
      } else {
        this.rtConfig = res.data
      }
    },
    async saveRt(rts) {
      this.selectPannel.reportTables = rts
      const res = await UpdatePannelRt({
        report_tables: rts.join(','),
        id: this.selectPannel.pannelId,
        appid: this.$store.state.baseData.EsConnectID
      })
      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })
        return
      } else {
        this.$message({
          type: 'success',
          offset: 60,
          message: res.msg
        })
      }
      this.GetPannelList()
      // this.loadReportTables()
    },
    compare(property) {
      return function(a, b) {
        var value1 = a[property]
        var value2 = b[property]
        return value2 - value1
      }
    },
    openAddPannelRt() {
      this.addPannelRtVisible = true
    },
    async loadReportTables() {
      this.echartBoxShow = false
      this.$nextTick(() =>
        this.echartBoxShow = true
      )
    },
    async deleteDir() {
      if (this.checkedKeys.length == 0) return

      this.$confirm('确定删除吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          for (const v of this.checkedKeys) {
            const res = await DeleteDir({ id: v, appid: this.$store.state.baseData.EsConnectID })
            if (res.code != 0) {
              this.$message({
                type: 'error',
                offset: 60,
                message: res.msg
              })
              return
            }
          }
          this.$message({
            type: 'success',
            offset: 60,
            message: '删除成功'
          })
          this.GetPannelList()
        })
        .catch(err => {
          console.error(err)
        })
    },
    async copyPannel() {
      const res = await CopyPannel(this.copyPannelForm)
      if (res.code != 0) {
        this.$message({
          offset: 60,
          type: 'error',
          message: res.msg
        })
        return
      }

      this.$message({
        type: 'success',
        offset: 60,
        message: res.msg
      })
      this.GetPannelList()
    },
    deletePannel(id, name) {
      this.$confirm('确定删除(' + name + ')该面板吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await DeletePannel({ appid: this.$store.state.baseData.EsConnectID, id: id })
          if (res.code != 0) {
            this.$message({
              type: 'error',
              offset: 60,
              message: res.msg
            })
            return
          }

          this.$message({
            type: 'success',
            offset: 60,
            message: res.msg
          })
          this.GetPannelList()
        })
        .catch(err => {
          console.error(err)
        })
    },
    dirClick({ key, domEvent }) {
    },
    openCopyPannel(report_tables, pannel_name, fid) {
      this.copyPannelForm.report_tables = report_tables.length == 0 ? '' : report_tables.join(',')
      this.copyPannelForm.pannel_name = pannel_name
      this.copyPannelForm.folder_id = fid
      this.copyPannelVisible = true
    },
    openMovePannelDialog(id, fid) {
      this.movePannelForm.id = id
      this.movePannelForm.folder_id = fid
      this.movePannelVisible = true
    },
    async movePannel() {
      const form = JSON.parse(JSON.stringify(this.movePannelForm))
      form['appid'] = this.$store.state.baseData.EsConnectID
      const res = await MovePannel2Dir(form)

      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })
        return
      }

      this.$message({
        type: 'success',
        offset: 60,
        message: res.msg
      })
      this.GetPannelList()

      this.movePannelVisible = false
    },
    async rename(id) {
      this.$prompt('请输入面板新名称', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }).then(async({ value }) => {
        const res = await Rename({ appid: this.$store.state.baseData.EsConnectID, id: id, pannel_name: value })
        if (res.code != 0) {
          this.$message({
            offset: 60,
            type: 'error',
            message: res.msg
          })
          return
        }

        this.$message({
          type: 'success',
          offset: 60,
          message: res.msg
        })
        this.GetPannelList()
      }).catch(err => {
        console.log(err)
      })
    },
    async addPannel() {
      if (this.addPannelForm.folder_id == '' || this.addPannelForm.folder_id == 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: '文件夹名不能为空'
        })
        return
      }
      if (this.addPannelForm.pannel_name == '') {
        this.$message({
          type: 'error',
          offset: 60,
          message: '看板名不能为空'
        })
        return
      }

      const form = JSON.parse(JSON.stringify(this.addPannelForm))
      form['appid'] = this.$store.state.baseData.EsConnectID
      const res = await NewPannel(form)

      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })
        return
      }

      this.$message({
        type: 'success',
        offset: 60,
        message: res.msg
      })
      this.GetPannelList()

      this.addPannelVisible = false
    },
    openAddPannel() {
      this.addPannelForm.folder_id = ''
      this.addPannelForm.pannel_name = ''
      this.addPannelVisible = true
    },
    async createDirReq(folder_name) {
      const res = await NewDir({
        appid: this.$store.state.baseData.EsConnectID,
        folder_name: folder_name
      })
      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })
        return
      } else {
        this.$message({
          type: 'success',
          offset: 60,
          message: res.msg
        })
        await this.GetPannelList()
      }
    },
    newDir() {
      this.$prompt('请输入文件夹名称', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消'
      }).then(({ value }) => {
        this.createDirReq(value)
      }).catch(err => {
        console.log(err)
      })
    },
    lookPannel(id, name, reportTables, canShare, managers, create_by) {
      this.selectPannel.pannelId = id
      this.selectPannel.pannelName = name
      this.selectPannel.reportTables = reportTables
      this.selectPannel.managers = managers.split(',')
      this.selectPannel.canShare = canShare
      this.selectPannel.create_by = create_by
    },
    async GetPannelList() {
      const res = await GetPannelList({ appid: this.$store.state.baseData.EsConnectID })
      if (res.code != 0) {
        this.$message({
          type: 'error',
          offset: 60,
          message: res.msg
        })
        return
      } else {
        if (res.data != null) {
          const tmp = {}
          res.data = filterData(res.data, this.panelStr.trim())

          let pannelId = 0
          if (this.LastSelectKey.length > 0) {
            const tmp = this.LastSelectKey[0].split('.')
            pannelId = Number(tmp[1])
          }
          this.openKeys = ['dashboard']
          for (const v of res.data) {
            if (pannelId != 0 && pannelId == v.pannel_id) {
              this.lookPannel(v.pannel_id, v['pannel_name'], v['report_tables'] == '' ? [] : v['report_tables'].split(','), v.folder_type == 0, v.managers, v.create_by)
            }

            this.openKeys.push(v['folder_id'])

            if (!tmp.hasOwnProperty(v['folder_id'])) {
              const obj = {
                folder_id: v['folder_id'],
                folder_name: v['folder_name'],
                folder_type: v['folder_type'],
                childList: []
              }

              if (v['pannel_id'] > 0) {
                obj.childList.push({
                  pannel_id: v['pannel_id'],
                  pannel_name: v['pannel_name'],
                  create_by: v['create_by'],
                  managers: v['managers'],
                  report_tables: v['report_tables'] == '' ? [] : v['report_tables'].split(',')
                })
              }
              tmp[v['folder_id']] = obj
            } else {
              if (v['pannel_id'] > 0) {
                tmp[v['folder_id']].childList.push(
                  {
                    pannel_id: v['pannel_id'],
                    pannel_name: v['pannel_name'],
                    create_by: v['create_by'],
                    managers: v['managers'],
                    report_tables: v['report_tables'] == '' ? [] : v['report_tables'].split(',')
                  }
                )
              }
            }
          }

          this.pannelList = []
          this.treeData = []

          Object.keys(tmp).forEach(key => {
            if (tmp[key].folder_type == 0) {
              this.treeData.push({ title: tmp[key].folder_name, key: tmp[key].folder_id })
            }
            this.pannelList.push(tmp[key])
          })
          this.pannelList = this.pannelList.sort(this.compare('folder_type'))
        }
      }
    },
    moment,
    filterDateCall(dates, dateStrings) {
      this.filterDate = dateStrings
    },
    cleanDateFiter() {
      this.filterDate = [
        moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD'),
        moment().startOf('day').subtract(1, 'days').format('YYYY-MM-DD')
      ]
    },
    onOpenChange(openKeys) { // 当菜单被展开时触发此处
      this.openKeys = openKeys
    }

  }

}
</script>

<style scoped src="@/styles/dashbord-index.css"/>
