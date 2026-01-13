import { computed } from 'vue'
import { useBreakpoints } from '@vueuse/core'

const breakpoints = useBreakpoints({
  sm: 640,
  md: 768,
  lg: 1024,
  xl: 1280
})

export function useResponsive() {
  const isMobile = breakpoints.smaller('md')
  const isDesktop = breakpoints.greaterOrEqual('md')

  const breakpoint = computed(() => {
    if (breakpoints.smaller('md').value) return 'sm'
    if (breakpoints.smaller('lg').value) return 'md'
    if (breakpoints.smaller('xl').value) return 'lg'
    return 'xl'
  })

  return { breakpoint, isMobile, isDesktop }
}

