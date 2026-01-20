// 表单配置类型定义
export interface FormSchema {
  field: string;
  label: string;
  component: 'Input' | 'Select' | 'DatePicker' | 'RangePicker' | 'InputNumber' | 'Textarea' | 'Switch' | 'TreeSelect' | 'RadioGroup' | 'IconPicker' | 'MonacoEditor';
  componentProps?: Record<string, any>;
  colProps?: {
    span?: number;
    offset?: number;
  };
  defaultValue?: any;
  required?: boolean;
  rules?: any[];
  show?: boolean | ((model: Record<string, any>) => boolean);
  helpMessage?: string;
}

export interface FormConfig {
  schemas: FormSchema[];
  labelWidth?: number | string;
  labelCol?: {
    span?: number;
  };
  wrapperCol?: {
    span?: number;
  };
}
