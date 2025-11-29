import React from 'react';
import { Button, Form, Input } from 'antd';
import axios from "@/axios/axios";
import Link from "next/link";
import router from "next/router";

const onFinish = (values: any) => {
    axios.post("/users/login", values)
        .then((res) => {
            if(res.status != 200) {
                alert(res.statusText);
                return
            }
            if(typeof res.data == 'string') {
                alert(res.data);
            } else {
                const msg = res.data?.msg || JSON.stringify(res.data)
                alert(msg);
                if(res.data.code == 0) {
                    router.push('/articles/list')
                }
            }
        }).catch((err) => {
            alert(err);
    })
};

const onFinishFailed = (errorInfo: any) => {
    alert("输入有误")
};

const LoginForm: React.FC = () => {
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
                    <h2 style={{ marginBottom: 8, fontSize: 28, color: '#111827' }}>欢迎回来</h2>
                    <p style={{ margin: 0, color: '#6b7280' }}>登录管理你的 Webook 内容</p>
                </div>
                <Form
                    layout="vertical"
                    name="login"
                    initialValues={{ remember: true }}
                    onFinish={onFinish}
                    onFinishFailed={onFinishFailed}
                    autoComplete="off"
                >
                    <Form.Item
                        label="邮箱"
                        name="email"
                        rules={[{ required: true, message: '请输入邮箱' }]}
                    >
                        <Input size="large" placeholder="your@email.com" />
                    </Form.Item>

                    <Form.Item
                        label="密码"
                        name="password"
                        rules={[{ required: true, message: '请输入密码' }]}
                    >
                        <Input.Password size="large" placeholder="请输入密码" />
                    </Form.Item>

                    <Form.Item style={{ marginBottom: 16 }}>
                        <Button
                            type="primary"
                            htmlType="submit"
                            size="large"
                            style={{ width: '100%' }}
                        >
                            登录
                        </Button>
                    </Form.Item>

                    <div style={{
                        display: 'flex',
                        flexDirection: 'column',
                        gap: 8,
                        textAlign: 'center'
                    }}>
                        <Link href={"/users/login_sms"}>使用手机号登录</Link>
                        <Link href={"/users/login_wechat"}>微信扫码登录</Link>
                        <Link href={"/users/signup"}>没有账号？立即注册</Link>
                    </div>
                </Form>
            </div>
        </div>
    )
};

export default LoginForm;