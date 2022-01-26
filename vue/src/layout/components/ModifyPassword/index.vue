<template>
  <el-dialog
    :close-on-click-modal="false"
    :visible.sync="dialogVisible"
    title="修改密码"
  >
    <el-form label-width="100px" label-position="left">
      <el-form-item label="密码">
        <el-input v-model="password" placeholder="密码" />
      </el-form-item>
    </el-form>
    <div style="text-align:right;">
      <el-button type="danger" icon="el-icon-close" @click="dialogVisible=false">取消</el-button>
      <el-button type="primary" icon="el-icon-check" @click="confirmRole">确认</el-button>
    </div>
  </el-dialog>
</template>

<script>
export default {
  name: 'Index',
  data() {
    return {
      password: ''
    }
  },
  methods: {

  }

}
</script>

<style scoped>

</style>
