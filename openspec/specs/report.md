# Spec: Report Generator

## Report Structure

1. 一、IV 期限结构概览
2. 二、IV vs RV 分析
3. 三、波动率偏斜（Skew）
4. 四、远期波动率矩阵（Forward IV）
5. 五、期权流向
6. 六、总结与交易观点

## Analysis Rules

### IV Change
- > 5%: 大幅上升/下降
- 2-5%: 小幅上升/下降
- < 2%: 基本持平

### VRP
- > 5: 期权溢价偏高
- 0-5: 定价合理
- < 0: 期权折价

### Skew
- > 1.05: 看跌偏斜
- 0.95-1.05: 中性
- < 0.95: 看涨偏斜

### Flow
- call_buys/put_buys > 1.5: 看涨
- put_buys/call_buys > 1.5: 看跌
- 否则: 均衡
