// 获取title
import defaultSettings from '@/settings'

const title = defaultSettings.title || '埋点数据分析中台'

export default function getPageTitle(pageTitle) {
  if (pageTitle) {
    return `${pageTitle} - ${title}`
  }
  return `${title}`
}
