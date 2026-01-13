<script setup lang="ts">
import { initUserSettingsFromStorage, useUserSettings } from '@/composables/useI18n'
import { onMounted } from 'vue'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

const { timezone, locale, setTimezone, setLocale } = useUserSettings()

onMounted(() => {
  initUserSettingsFromStorage()
})
</script>

<template>
  <div class="mx-auto w-full max-w-3xl px-4 py-4">
    <div class="text-lg font-semibold">设置</div>

    <div class="mt-4 grid gap-4">
      <div>
        <div class="mb-1 text-sm text-muted-foreground">时区</div>
        <Input :model-value="timezone" @update:model-value="setTimezone(String($event))" />
      </div>

      <div>
        <div class="mb-1 text-sm text-muted-foreground">语言</div>
        <div class="flex gap-2">
          <Button variant="outline" :class="locale === 'zh' ? 'border-primary' : ''" @click="setLocale('zh')">
            中文
          </Button>
          <Button variant="outline" :class="locale === 'en' ? 'border-primary' : ''" @click="setLocale('en')">
            English
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>

