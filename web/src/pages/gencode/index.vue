<template>

  <el-form
      label-position="left"
      style="max-width: 600px;"
  >
    <el-form-item label="选择表">
      <el-select v-model="selectTableRef">
        <el-option :key="v.table_name" :label="v.table_comment+' ' + v.table_name" :value="v.table_name"
                   v-for="v in tablesRef"/>
      </el-select>
    </el-form-item>
  </el-form>

  <el-button @click="loadTableColumns">确定</el-button>

  <el-card style="margin-top: 15px">

    <el-form :inline="true" label-position="right" label-width="100px">
      <el-form-item label="包名称">

        <el-input v-model="columnsRef.module">
          <template #suffix>
            <el-tooltip>
              <el-icon>
                <QuestionFilled/>
              </el-icon>
              <template #content>
                <span>go包模块位置,为空标识不分模块</span>
              </template>
            </el-tooltip>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item label="文件名称">
        <el-input v-model="columnsRef.file_name"/>
      </el-form-item>
      <el-form-item label="struct名称">
        <el-input v-model="columnsRef.struct_name"/>
      </el-form-item>
      <el-form-item label="中文名称">
        <el-input v-model="columnsRef.ch_name"/>
      </el-form-item>
    </el-form>
  </el-card>
  <el-divider/>


  <el-table :data="columnsRef?.fields" :border="true">
    <el-table-column prop="field_name" label="字段名称" width="140">
      <template #default="scope">
        <el-input type="text" v-model="scope.row.field_name"/>
      </template>

    </el-table-column>
    <el-table-column prop="field_zh_name" label="中文名称" width="120">
      <template #default="scope">
        <el-input type="text" v-model="scope.row.field_zh_name"/>
      </template>
    </el-table-column>
    <el-table-column prop="field_type" label="数据类型" width="180">
      <template #default="scope">
        <el-input type="text" v-model="scope.row.field_type"/>
      </template>

    </el-table-column>
    <el-table-column prop="field_comment" label="注释" width="130">
      <template #default="scope">
        <el-input type="text" v-model="scope.row.field_comment"/>
      </template>
    </el-table-column>
    <el-table-column prop="field_json" label="数据库字段" min-width="120">
      <template #default="scope">
        <el-input type="text" v-model="scope.row.field_json"/>
      </template>
    </el-table-column>
    <el-table-column prop="show" label="是否显示">
      <template #default="scope">
        <input type="checkbox" v-model="scope.row.show"/>
      </template>
    </el-table-column>
    <el-table-column prop="require" label="编辑是否必填">
      <template #default="scope">
        <input type="checkbox" v-model="scope.row.require"/>
      </template>
    </el-table-column>

    <el-table-column prop="search_type" label="搜索条件" min-width="100">
      <template #default="scope">
        <el-select v-model="scope.row.search_type" clearable>
          <el-option v-for="s in searchOpt" :label="s" :value="s"/>
        </el-select>
      </template>
    </el-table-column>

  </el-table>

  <el-dialog v-model="showCodeRef">
    <el-tabs type="border-card">
      <el-tab-pane v-for="(v,k) in previewCodeRef?.codes" :label="v.file_name">
        <hljsVuePlugin.component :code="v.code"/>
      </el-tab-pane>

    </el-tabs>


  </el-dialog>
  <el-button @click="preview">预览</el-button>
  <el-button @click="download">下载</el-button>
</template>

<script setup lang="ts">
import {downloadCode, getTableColumns, getTables, Preview, previewCode, PreviewRes, TableRes} from '@/api/module/auto';
import {reactive, ref} from 'vue';
import 'highlight.js/styles/stackoverflow-light.css'
import 'highlight.js/lib/common';

import hljsVuePlugin from "@highlightjs/vue-plugin";
import {ElNotification} from "element-plus";


const searchOpt = ["=", "like", "between"]

const tablesRef = ref<TableRes[]>()
const selectTableRef = ref('')
const columnsRef = reactive<Preview>({
  router_path: "",
  go_out_dir: "", js_out_dir: "", package_prefix: "",
  table_name: '',
  struct_name: '',
  fields: [],
  file_name: '',
  comment: '',
  module: '',
  ch_name: ''
})
const showCodeRef = ref(false)
const previewCodeRef = ref<PreviewRes>()

async function initData() {
  tablesRef.value = await getTables()
  console.log(tablesRef.value);

}

initData()

async function loadTableColumns() {
  const tableColumns = await getTableColumns(selectTableRef.value)
  columnsRef.table_name = tableColumns.table_name
  columnsRef.fields = tableColumns.fields
  columnsRef.module = tableColumns.module
  columnsRef.struct_name = tableColumns.struct_name
  columnsRef.file_name = tableColumns.file_name
  columnsRef.comment = tableColumns.comment
  columnsRef.ch_name = tableColumns.ch_name
}

async function preview() {

  columnsRef.go_out_dir = localStorage.getItem(`go_target`)!!
  columnsRef.package_prefix = localStorage.getItem(`package_prefix`)!!
  columnsRef.js_out_dir =  localStorage.getItem(`web_target`)!!

  const code = await previewCode(columnsRef)

  console.log(code)
  showCodeRef.value = !showCodeRef.value;
  previewCodeRef.value = code

}

async function download() {
  try {
    columnsRef.go_out_dir = localStorage.getItem(`go_target`)!!
    columnsRef.package_prefix = localStorage.getItem(`package_prefix`)!!
    columnsRef.js_out_dir =  localStorage.getItem(`web_target`)!!
    await downloadCode(columnsRef)
  } catch (e) {

  }
  ElNotification({
    title: '通知',
    message: "下载成功!",
    duration: 2000,
  })


}

</script>


<style scoped></style>
