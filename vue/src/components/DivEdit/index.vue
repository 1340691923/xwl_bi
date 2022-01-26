<template>
  <div
    contenteditable="plaintext-only"
    @input="changeDivText($event)"
    @focus="divFocus"
    @blur="divBlur"
    v-html="value"
  >{{ value }}
  </div>
</template>

<script>
export default {
  name: 'Index',
  props: {
    value: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      editText: ''
    }
  },
  created() {
    this.editText = this.value
  },
  methods: {
    // 文本改变时事件
    changeDivText(event) {
      this.editText = event.target.innerHTML
    },
    divFocus() {
      this.editText = this.value
    },
    divBlur() {
      this.$emit('input', this.editText)
    }
  }
}
</script>
