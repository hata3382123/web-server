import React from 'react';
import { Button, Card, Form, Input, Typography } from 'antd';
import axios from "@/axios/axios";
import Link from "next/link";
import router from "next/router";

const { Title, Text } = Typography;

const onFinish = (values: any) => {
    axios.post("/users/signup", values)
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
                    router.push('/users/login')
                }
            }

        }).catch((err) => {
            alert(err);
    })
};

const onFinishFailed = (errorInfo: any) => {
    alert("输入有误")
};

const SignupForm: React.FC = () => (
    <div
        style={{
            minHeight: '100vh',
            background: 'linear-gradient(135deg, #e0f2ff 0%, #f5f7ff 100%)',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            padding: '40px 16px',
        }}
    >
        <Card
            style={{ width: '100%', maxWidth: 420, boxShadow: '0 12px 40px rgba(15, 71, 147, 0.12)' }}
            bordered={false}
        >
            <div style={{ textAlign: 'center', marginBottom: 32 }}>
                <Title level={3} style={{ marginBottom: 8 }}>
                    创建您的账号
                </Title>
                <Text type="secondary">加入 Webook，探索更多精彩内容</Text>
            </div>
            <Form
                layout="vertical"
                name="signup"
                onFinish={onFinish}
                onFinishFailed={onFinishFailed}
                autoComplete="off"
            >
                <Form.Item
                    label="邮箱"
                    name="email"
                    rules={[
                        { required: true, message: '请输入邮箱' },
                        { type: 'email', message: '邮箱格式不正确' },
                    ]}
                >
                    <Input size="large" placeholder="name@example.com" />
                </Form.Item>

                <Form.Item
                    label="密码"
                    name="password"
                    rules={[
                        { required: true, message: '请输入密码' },
                        { min: 6, message: '密码至少 6 位' },
                    ]}
                >
                    <Input.Password size="large" placeholder="至少 6 位，区分大小写" />
                </Form.Item>

                <Form.Item
                    label="确认密码"
                    name="confirmPassword"
                    dependencies={['password']}
                    rules={[
                        { required: true, message: '请确认密码' },
                        ({ getFieldValue }) => ({
                            validator(_, value) {
                                if (!value || getFieldValue('password') === value) {
                                    return Promise.resolve();
                                }
                                return Promise.reject(new Error('两次输入的密码不一致'));
                            },
                        }),
                    ]}
                >
                    <Input.Password size="large" placeholder="再次输入密码" />
                </Form.Item>
                <Form.Item style={{ marginTop: 32 }}>
                    <Button type="primary" htmlType="submit" size="large" block>
                        注册
                    </Button>
                    <div style={{ textAlign: 'center', marginTop: 16 }}>
                        <Text type="secondary">
                            已有账号？ <Link href={"/users/login"}>立即登录</Link>
                        </Text>
                    </div>
                </Form.Item>
            </Form>
        </Card>
    </div>
);

export default SignupForm;