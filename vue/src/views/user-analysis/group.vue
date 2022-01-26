<template>
  <div class="app-container">
    <el-card class="box-card" style="min-height:  850px;">
      <div
        style="height: 50px;line-height: 50px;display: flex;align-items: center;justify-content: space-between;border-bottom: 1px solid #f0f2f5"
      >

        <a-tooltip placement="right" style="cursor: pointer">
          <template slot="title">
            <span>将具有共同特征的用户组建成分群，方便在各种模型中利用分群进行细分筛选分析</span>
          </template>
          <span class="title_xwl" style="color: #202d3f">&nbsp;&nbsp;用户分群 <a-icon
            type="question-circle"
          />
          </span>
        </a-tooltip>
        <div>
          <a-input-search v-model="input" placeholder="请输入搜索" style="width: 200px" />
        </div>
      </div>
      <page-table
        v-if="tableShow"
        :border="false"
        :connect-loading="loading"
        :table-list="tableData"
        :table-info="tableInfo"
        :input="input"
      >
        <el-table-column
          slot="operate"
          align="center"
          prop="id"
          label="分群ID"
          width="100"
        />
        <el-table-column
          slot="operate"
          align="center"
          label="分群名"
        >
          <template slot-scope="scope">
            <div v-if="scope.row.isEdit">
              <el-input v-model="scope.row.group_name" size="small" />
            </div>
            <div v-else>
              {{ scope.row.group_name }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          slot="operate"
          align="center"
          label="分群备注"
        >
          <template slot-scope="scope">
            <div v-if="scope.row.isEdit">
              <el-input v-model="scope.row.group_remark" size="small" />
            </div>
            <div v-else>
              {{ scope.row.group_remark }}
            </div>
          </template>
        </el-table-column>
        <el-table-column
          slot="operate"
          align="center"
          sortable
          prop="create_time"
          label="创建时间"
          width="150"
        />
        <el-table-column
          slot="operate"
          align="center"
          sortable
          prop="update_time"
          label="更新时间"
          width="150"
        />
        <el-table-column
          slot="operate"
          align="center"
          sortable
          prop="user_count"
          label="用户群人数"
          width="100"
        />

        <el-table-column
          slot="operate"
          align="center"
          fixed="right"
          label="操作"
          width="500"
        >
          <template slot-scope="scope">
            <el-button
              v-if="scope.row.isEdit"
              size="mini"
              icon="el-icon-check"
              @click.native.prevent="Save(scope.row.index,scope.row)"
            >
              保存
            </el-button>
            <el-button
              v-if="!scope.row.isEdit"
              type="primary"
              size="mini"
              icon="el-icon-edit"
              @click.native.prevent="Update(scope.row.index)"
            >
              修改
            </el-button>
            <el-button
              type="danger"
              size="mini"
              icon="el-icon-close"
              @click.native.prevent="Delete(scope.row.id,scope.$index)"
            >
              删除
            </el-button>
            <el-button
              type="success"
              size="mini"
              icon="el-icon-user"
              @click.native.prevent="drillDown(scope.row.user_list)"
            >
              进入该用户分群
            </el-button>
          </template>
        </el-table-column>
      </page-table>

    </el-card>
  </div>
</template>

<script>
import { DeleteUserGroup, ModifyUserGroup, UserGroupList } from '@/api/user-group'

export default {
  name: 'Tag',
  components: {
    'PageTable': () => import('@/components/PageTable')
  },
  data() {
    return {
      tableInfo: [{ slot: 'operate' }],
      tableShow: true,
      input: '',
      loading: false,
      tableData: []
    }
  },
  mounted() {
    this.init()
  },
  methods: {
    drillDown(ui) {
      this.$store.dispatch('baseData/SETUI', ui)
      this.$router.push({ path: '/user-analysis/user_list' })
    },
    async Delete(id, index) {
      this.$confirm('确定删除该分群吗?', '警告', {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const form = { 'appid': this.$store.state.baseData.EsConnectID }
          form['id'] = id
          const res = await DeleteUserGroup(form)
          if (res.code != 0) {
            this.$message({
              showClose: true,
              offset: 60,
              type: 'error',
              message: res.msg
            })
            return
          }
          this.$message({
            showClose: true,
            offset: 60,
            type: 'success',
            message: res.msg
          })
          this.tableData.splice(index, 1)
          this.refreshTable()
        })
        .catch(err => {
          console.error(err)
        })
    },
    async Save(index, row) {
      const form = { 'appid': this.$store.state.baseData.EsConnectID }
      form['id'] = row['id']
      form['name'] = row['group_name']
      form['remark'] = row['group_remark']

      const res = await ModifyUserGroup(form)
      if (res.code != 0) {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'error',
          message: res.msg
        })
        return
      }
      this.$message({
        showClose: true,
        offset: 60,
        type: 'success',
        message: res.msg
      })

      for (const k in this.tableData) {
        if (this.tableData[k].index == index) {
          this.tableData[k].isEdit = false
          break
        }
      }
      this.refreshTable()
    },
    Update(index) {
      for (const k in this.tableData) {
        if (this.tableData[k].index == index) {
          this.tableData[k].isEdit = true
          break
        }
      }
      this.refreshTable()
    },
    refreshTable() {
      this.tableShow = false
      this.$nextTick(() => {
        this.tableShow = true
      })
    },
    async init() {
      this.tableData = []
      const form = { 'appid': this.$store.state.baseData.EsConnectID }
      this.loading = true
      const res = await UserGroupList(form)
      this.loading = false
      if (res.code != 0) {
        this.$message({
          showClose: true,
          offset: 60,
          type: 'error',
          message: res.msg
        })
        this.refreshTable()
        return
      }
      this.$message({
        showClose: true,
        offset: 60,
        type: 'success',
        message: res.msg
      })
      if (res.data == null) {
        res.data = []
      }

      this.tableData = res.data

      for (const i in this.tableData) {
        this.tableData[i].index = i
        this.tableData[i].isEdit = false
      }

      this.refreshTable()
    }
  }
}
</script>

<style scoped>

</style>
