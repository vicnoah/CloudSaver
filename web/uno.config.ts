import { defineConfig, presetUno, presetAttributify, presetIcons } from 'unocss'

export default defineConfig({
  presets: [
    presetUno(),
    presetAttributify(),
    presetIcons({
      collections: {
        carbon: () => import('@iconify-json/carbon/icons.json').then(i => i.default),
        mdi: () => import('@iconify-json/mdi/icons.json').then(i => i.default),
      }
    }),
  ],
  shortcuts: {
    'btn': 'px-4 py-2 rounded cursor-pointer transition-colors',
    'btn-primary': 'btn bg-blue-500 text-white hover:bg-blue-600',
    'btn-secondary': 'btn bg-gray-200 text-gray-800 hover:bg-gray-300',
    'input-base': 'px-3 py-2 border border-gray-300 rounded focus:outline-none focus:border-blue-500',
    'card': 'bg-white rounded-lg shadow p-4',
  },
  theme: {
    colors: {
      primary: '#409EFF',
      success: '#67C23A',
      warning: '#E6A23C',
      danger: '#F56C6C',
      info: '#909399',
    },
    breakpoints: {
      'sm': '640px',
      'md': '768px',
      'lg': '1024px',
      'xl': '1280px',
    },
  },
})
