// plugins/echarts.ts
import * as echarts from 'echarts/core';
import {
  LineChart,
  PieChart
} from 'echarts/charts';
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent
} from 'echarts/components';
import {
  CanvasRenderer
} from 'echarts/renderers';

echarts.use([
  LineChart,
  PieChart,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
  CanvasRenderer
]);

export default defineNuxtPlugin(() => {
  return {
    provide: {
      echarts
    }
  };
});
