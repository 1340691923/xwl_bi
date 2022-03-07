<template>
  <div>
    <el-button type="primary" @click="openDialog">保存报表</el-button>
    <el-button type="primary" plain @click="go">计算</el-button>
    <el-dialog :close-on-click-modal="false" :visible.sync="dialogVisible" title="添加/更新报表" @close="formClose">
      <el-form :model="form" label-width="120px" label-position="left">
        <el-form-item label="报表名称:">
          <el-input v-model="form.name" style="width: 70%" placeholder="相同报表名称的报表将会被覆盖更新" />
          <el-button type="primary" @click="findNameCount()">检测是否有同名</el-button>
        </el-form-item>
        <el-form-item label="报表备注:">
          <el-input v-model="form.remark" type="textarea" :autosize="{ minRows: 4, maxRows: 8}" placeholder="报表备注" />
        </el-form-item>
      </el-form>
      <div style="text-align:right;">
        <el-button type="danger" icon="el-icon-close" @click="formClose">返回</el-button>
        <el-button type="primary" icon="el-icon-plus" @click="addForm">添加/更新</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { AddReportTable, FindNameCount } from '@/api/pannel'

export default {
  name: 'AddReportTable',
  props: {
    rtType: {
      type: Number,
      default: 0
    },
    name: {
      type: String,
      default: ''
    },
    remark: {
      type: String,
      default: ''
    },
    data: {
      type: Object,
      default: {}
    }
  },
  data() {
    return {
      dialogVisible: false,
      form: {
        name: this.Name,
        remark: this.Remark,
        rt_type: this.rtType,
        appid: this.$store.state.baseData.EsConnectID
      }
    }
  },
  watch: {
    Name(oldV, newV) {
      this.form.name = this.Name
      this.form.remark = this.Remark
    }
  },
  methods: {
    async findNameCount() {
      const res = await FindNameCount(this.form)
      if (res.code != 0) {
        this.$message({
          offset: 60,

          type: 'error',
          message: res.msg
        })
        return
      } else {
        if (res.data > 0) {
          this.$notify({
            title: 'Success',
            dangerouslyUseHTMLString: true,
            message: `
                  <div>已经存在该报表，再次操作将更新这张报表</div>
                `,
            type: 'warning'
          })
        } else {
          this.$notify({
            title: 'Success',
            dangerouslyUseHTMLString: true,
            message: `
                  <div>不存在该报表，再次操作将新增</div>
                `,
            type: 'success'
          })
        }
      }
    },
    async addForm() {
      const input = JSON.parse(JSON.stringify(this.form))
      input['data'] = JSON.stringify(this.data)
      const res = await AddReportTable(input)
      if (res.code != 0) {
        this.$message({
          offset: 60,

          type: 'error',
          message: res.msg
        })
        return
      } else {
        this.$message({
          offset: 60,

          type: 'success',
          message: res.msg
        })
        this.dialogVisible = false
      }
    },
    formClose() {
      this.dialogVisible = false
    },
    openDialog() {
      this.dialogVisible = true
    },
    go() {
      this.$emit('go')
    }
  }
}
</script>

<style scoped>

</style>
