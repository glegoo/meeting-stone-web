import { computed, ref } from 'vue'
import { formatInTimeZone } from 'date-fns-tz'

const timezone = ref<string>(Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC')
const locale = ref<'zh' | 'en'>('zh')

export function useUserSettings() {
  const setTimezone = (tz: string) => {
    timezone.value = tz
    localStorage.setItem('timezone', tz)
  }
  const setLocale = (lc: 'zh' | 'en') => {
    locale.value = lc
    localStorage.setItem('locale', lc)
  }

  return { timezone, locale, setTimezone, setLocale }
}

export function initUserSettingsFromStorage() {
  const tz = localStorage.getItem('timezone')
  const lc = localStorage.getItem('locale') as 'zh' | 'en' | null
  if (tz) timezone.value = tz
  if (lc) locale.value = lc
}

export function useActivityTime(deadlineUtc: string) {
  const settings = useUserSettings()

  return computed(() => {
    const date = new Date(deadlineUtc)
    // 先用固定格式；后续接入 vue-i18n 后再按语言切换
    return formatInTimeZone(date, settings.timezone.value, 'yyyy-MM-dd HH:mm')
  })
}

