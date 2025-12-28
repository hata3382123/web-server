import React, { useState, useEffect } from 'react';
import { Button, Form, Input, message, Space } from 'antd';
import axios from "@/axios/axios";
import Link from "next/link";
import router from "next/router";

const LoginFormSMS: React.FC = () => {
    const [form] = Form.useForm();
    const [countdown, setCountdown] = useState(0);
    const [loading, setLoading] = useState(false);
    const [sending, setSending] = useState(false);

    // 倒计时效果
    useEffect(() => {
        if (countdown > 0) {
            const timer = setTimeout(() => {
                setCountdown(countdown - 1);
            }, 1000);
            return () => clearTimeout(timer);
        }
    }, [countdown]);

    // 手机号格式验证
    const validatePhone = (_: any, value: string) => {
        if (!value) {
            return Promise.reject(new Error('请输入手机号码'));
        }
        const phoneRegex = /^1[3-9]\d{9}$/;
        if (!phoneRegex.test(value)) {
            return Promise.reject(new Error('请输入正确的手机号码'));
        }
        return Promise.resolve();
    };

    // 验证码格式验证
    const validateCode = (_: any, value: string) => {
        if (!value) {
            return Promise.reject(new Error('请输入验证码'));
        }
        if (!/^\d{6}$/.test(value)) {
            return Promise.reject(new Error('验证码为6位数字'));
        }
        return Promise.resolve();
    };

    const sendCode = async () => {
        try {
            const phone = form.getFieldValue("phone");
            if (!phone) {
                message.error('请先输入手机号码');
                return;
            }
            
            // 验证手机号格式
            try {
                await form.validateFields(['phone']);
            } catch {
                return;
            }

            setSending(true);
            const res = await axios.post("/users/login_sms/code/send", { phone });
            
            if (res.status === 200) {
                const msg = res?.data?.msg || "验证码已发送";
                message.success(msg);
                setCountdown(60); // 开始60秒倒计时
            } else {
                message.error(res.statusText || "发送失败，请重试");
            }
        } catch (err: any) {
            const errorMsg = err?.response?.data?.msg || err?.message || "系统错误，请重试";
            message.error(errorMsg);
        } finally {
            setSending(false);
        }
    };

    const onFinish = async (values: any) => {
        setLoading(true);
        try {
            const res = await axios.post("/users/login_sms", values);
            
            if (res.status === 200) {
                if (res.data.code === 0) {
                    message.success('登录成功');
                    router.push('/users/profile');
                    return;
                }
                message.warning(res.data.msg || '登录失败');
            } else {
                message.error(res.statusText || '登录失败');
            }
        } catch (err: any) {
            const errorMsg = err?.response?.data?.msg || err?.message || "系统错误，请重试";
            message.error(errorMsg);
        } finally {
            setLoading(false);
        }
    };

    const onFinishFailed = () => {
        message.error('请检查输入信息');
    };

    return (
        <div style={{
            minHeight: '100vh',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '40px',
            background: 'linear-gradient(135deg, #f6f8ff 0%, #e0e7ff 50%, #fef3c7 100%)'
        }}>
            <div style={{
                width: '100%',
                maxWidth: 420,
                background: '#fff',
                borderRadius: 16,
                boxShadow: '0 20px 45px rgba(15, 23, 42, 0.15)',
                padding: '48px 40px'
            }}>
                <div style={{ marginBottom: 32, textAlign: 'center' }}>
                    <h2 style={{ marginBottom: 8, fontSize: 28, color: '#111827' }}>手机号登录</h2>
                    <p style={{ margin: 0, color: '#6b7280' }}>使用手机验证码快速登录或注册</p>
                </div>
                
                <Form
                    layout="vertical"
                    name="login_sms"
                    form={form}
                    initialValues={{ remember: true }}
                    onFinish={onFinish}
                    onFinishFailed={onFinishFailed}
                    autoComplete="off"
                >
                    <Form.Item
                        label="手机号码"
                        name="phone"
                        rules={[{ validator: validatePhone }]}
                    >
                        <Input 
                            size="large" 
                            placeholder="请输入11位手机号码"
                            maxLength={11}
                            style={{ fontSize: 16 }}
                        />
                    </Form.Item>

                    <Form.Item
                        label="验证码"
                        name="code"
                        rules={[{ validator: validateCode }]}
                    >
                        <Space.Compact style={{ width: '100%' }}>
                            <Input 
                                size="large" 
                                placeholder="请输入6位验证码"
                                maxLength={6}
                                style={{ fontSize: 16 }}
                            />
                            <Button
                                size="large"
                                type={countdown > 0 ? "default" : "primary"}
                                disabled={countdown > 0 || sending}
                                loading={sending}
                                onClick={sendCode}
                                style={{ 
                                    minWidth: 120,
                                    fontWeight: 500
                                }}
                            >
                                {countdown > 0 ? `${countdown}秒后重发` : '发送验证码'}
                            </Button>
                        </Space.Compact>
                    </Form.Item>

                    <Form.Item style={{ marginBottom: 16, marginTop: 8 }}>
                        <Button
                            type="primary"
                            htmlType="submit"
                            size="large"
                            loading={loading}
                            style={{ 
                                width: '100%',
                                height: 44,
                                fontSize: 16,
                                fontWeight: 500
                            }}
                        >
                            {loading ? '登录中...' : '登录/注册'}
                        </Button>
                    </Form.Item>

                    <div style={{
                        display: 'flex',
                        flexDirection: 'column',
                        gap: 8,
                        textAlign: 'center'
                    }}>
                        <Link href={"/users/login"} style={{ color: '#6b7280', textDecoration: 'none' }}>
                            使用邮箱登录
                        </Link>
                        <Link href={"/users/signup"} style={{ color: '#6b7280', textDecoration: 'none' }}>
                            没有账号？立即注册
                        </Link>
                    </div>
                </Form>
            </div>
        </div>
    );
};

export default LoginFormSMS;