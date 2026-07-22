-- +goose Up
DROP INDEX IF EXISTS public.idx_s_config_code;

ALTER TABLE public.s_config
  ALTER COLUMN data SET NOT NULL;

CREATE UNIQUE INDEX idx_s_config_code ON public.s_config USING btree (code);

INSERT INTO public.s_config (
  id, name, code, data, remark, create_by, created_time, update_by, updated_time
) VALUES
  (
    1880159541355580001,
    '微信集成配置',
    'integration.wechat',
    '{"enabled":false,"appId":"","secret":"","templateId":""}'::jsonb,
    '微信小程序登录与消息能力运行期配置',
    0, NOW(), 0, NOW()
  ),
  (
    1880159541355580002,
    '短信集成配置',
    'integration.sms',
    '{"enabled":false,"accessKeyId":"","accessKeySecret":"","signName":"","templateCode":""}'::jsonb,
    '短信服务运行期配置',
    0, NOW(), 0, NOW()
  ),
  (
    1880159541355580003,
    '邮件集成配置',
    'integration.email',
    '{"enabled":false,"host":"","port":0,"username":"","password":"","from":""}'::jsonb,
    '邮件服务运行期配置',
    0, NOW(), 0, NOW()
  ),
  (
    1880159541355580004,
    '验证码配置',
    'auth.captcha',
    '{"image":{"enabled":false,"length":4,"width":120,"height":40,"expire":300},"sms":{"enabled":false,"length":6,"expire":300,"template":"SMS_CODE_TEMPLATE","provider":"aliyun"},"email":{"enabled":false,"length":6,"expire":300,"template":"验证码：%s"}}'::jsonb,
    '图形、短信和邮件验证码运行期配置',
    0, NOW(), 0, NOW()
  ),
  (
    1880159541355580005,
    '调度器配置',
    'scheduler',
    '{"enabled":true,"refreshIntervalSeconds":5,"jobs":{"data-cleanup":{"enabled":true,"cron":"0 0 2 * * *"}}}'::jsonb,
    '仅允许调度代码注册的后台任务',
    0, NOW(), 0, NOW()
  )
ON CONFLICT (code) DO UPDATE SET
  name = EXCLUDED.name,
  data = EXCLUDED.data,
  remark = EXCLUDED.remark,
  update_by = EXCLUDED.update_by,
  updated_time = NOW();

-- +goose Down
DELETE FROM public.s_config
WHERE code IN ('integration.wechat', 'integration.sms', 'integration.email', 'auth.captcha', 'scheduler');

DROP INDEX IF EXISTS public.idx_s_config_code;
CREATE INDEX idx_s_config_code ON public.s_config USING btree (code);

ALTER TABLE public.s_config
  ALTER COLUMN data DROP NOT NULL;
