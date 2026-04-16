<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { Form, Input, InputNumber, Button, Alert, Result } from 'ant-design-vue'
import { DatabaseOutlined, CheckCircleOutlined, LoadingOutlined } from '@ant-design/icons-vue'
import { apiTestConnection, apiSetup } from '@/api/init'

const router = useRouter()
let needsInit = true
const formRef = ref()

const form = reactive({
  host: '127.0.0.1',
  port: 3306,
  user: 'root',
  password: '',
  database: 'itfeeds',
})

const loading = ref(false)
const testing = ref(false)
const step = ref(0) // 0=填写, 1=测试中, 2=初始化中, 3=完成
const testResult = ref(null)
const errorMsg = ref('')

const rules = {
  host: [{ required: true, message: '请输入数据库主机' }],
  port: [{ required: true, message: '请输入端口' }],
  user: [{ required: true, message: '请输入用户名' }],
  database: [{ required: true, message: '请输入数据库名' }],
}

async function handleTest() {
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  testing.value = true
  testResult.value = null
  errorMsg.value = ''

  try {
    const res = await apiTestConnection(form)
    testResult.value = res
  } catch (e) {
    testResult.value = { success: false, error: e.message || '连接失败' }
  } finally {
    testing.value = false
  }
}

async function handleSetup() {
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  step.value = 2
  errorMsg.value = ''

  try {
    await apiSetup(form)
    needsInit = false
    step.value = 3
    setTimeout(() => {
      window.location.href = '/'
    }, 1500)
  } catch (e) {
    step.value = 0
    errorMsg.value = e.message || '初始化失败'
  }
}
</script>

<template>
  <div class="init-page">
    <div class="init-container">
      <!-- Logo -->
      <div class="init-logo">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" width="28" height="28" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M4 11a9 9 0 0 1 9 9"/>
            <path d="M4 4a16 16 0 0 1 16 16"/>
            <circle cx="5" cy="19" r="1"/>
          </svg>
        </div>
        <h1 class="logo-text">IT<span>Feeds</span></h1>
      </div>

      <p class="init-desc">首次使用需要配置 MySQL 数据库连接</p>

      <!-- 完成状态 -->
      <Result
        v-if="step === 3"
        status="success"
        title="初始化完成"
        sub-title="即将跳转到首页..."
      />

      <!-- 表单 -->
      <div v-else class="init-card">
        <Form ref="formRef" :model="form" :rules="rules" layout="vertical">
          <div class="form-row">
            <Form.Item label="数据库主机" name="host" class="form-item flex-1">
              <Input v-model:value="form.host" placeholder="127.0.0.1" :disabled="step >= 2" />
            </Form.Item>
            <Form.Item label="端口" name="port" class="form-item port-item">
              <InputNumber v-model:value="form.port" :min="1" :max="65535" style="width: 100%" :disabled="step >= 2" />
            </Form.Item>
          </div>
          <div class="form-row">
            <Form.Item label="用户名" name="user" class="form-item flex-1">
              <Input v-model:value="form.user" placeholder="root" :disabled="step >= 2" />
            </Form.Item>
            <Form.Item label="密码" name="password" class="form-item flex-1">
              <Input.Password v-model:value="form.password" placeholder="数据库密码" :disabled="step >= 2" />
            </Form.Item>
          </div>
          <Form.Item label="数据库名" name="database">
            <Input v-model:value="form.database" placeholder="itfeeds" :disabled="step >= 2" />
          </Form.Item>

          <!-- 测试结果 -->
          <Alert
            v-if="testResult"
            :type="testResult.success ? 'success' : 'error'"
            :message="testResult.success ? `连接成功 (MySQL ${testResult.version})` : testResult.error"
            show-icon
            style="margin-bottom: 16px"
          />

          <!-- 错误提示 -->
          <Alert
            v-if="errorMsg"
            type="error"
            :message="errorMsg"
            show-icon
            style="margin-bottom: 16px"
          />

          <!-- 操作按钮 -->
          <div class="form-actions">
            <Button
              v-if="step < 2"
              :loading="testing"
              @click="handleTest"
            >
              <template #icon><DatabaseOutlined /></template>
              测试连接
            </Button>
            <Button
              type="primary"
              :loading="step === 2"
              :disabled="!testResult?.success"
              @click="handleSetup"
            >
              <template #icon>
                <LoadingOutlined v-if="step === 2" />
                <CheckCircleOutlined v-else />
              </template>
              {{ step === 2 ? '初始化中...' : '开始初始化' }}
            </Button>
          </div>
        </Form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.init-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  padding: 24px;
}

.init-container {
  width: 100%;
  max-width: 460px;
}

.init-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  justify-content: center;
  margin-bottom: 8px;
}

.logo-icon {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, #d32f2f 0%, #f44336 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}

.logo-text {
  font-size: 28px;
  font-weight: 700;
  color: rgba(0, 0, 0, 0.85);
  letter-spacing: -0.5px;
  margin: 0;
}

.logo-text span {
  color: #d32f2f;
}

.init-desc {
  text-align: center;
  color: rgba(0, 0, 0, 0.45);
  margin-bottom: 32px;
  font-size: 14px;
}

.init-card {
  background: #fff;
  border-radius: 12px;
  padding: 32px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.form-row {
  display: flex;
  gap: 16px;
}

.flex-1 {
  flex: 1;
}

.port-item {
  width: 120px;
  flex: none;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

@media (max-width: 480px) {
  .init-card {
    padding: 24px 20px;
  }

  .form-row {
    flex-direction: column;
    gap: 0;
  }

  .port-item {
    width: 100%;
  }
}
</style>
