import { ProDescriptions } from '@ant-design/pro-components';
import React, { useState, useEffect } from 'react';
import { Button, Card, Avatar, Typography, Space } from 'antd';
import { UserOutlined, MailOutlined, PhoneOutlined } from '@ant-design/icons';
import axios from "@/axios/axios";

function Page() {
    let p: Profile = {Email: "", Phone: "", Nickname: "", Birthday:"", AboutMe: ""}
    const [data, setData] = useState<Profile>(p)
    const [isLoading, setLoading] = useState(false)

    useEffect(() => {
        setLoading(true)
        axios.get('/users/profile')
            .then((res) => res.data.data)
            .then((data) => {
                setData(data)
                setLoading(false)
            })
            .catch(() => {
                setLoading(false)
            })
    }, [])

    if (isLoading) return <p style={{ textAlign: 'center', marginTop: 80 }}>加载中...</p>
    if (!data) return <p style={{ textAlign: 'center', marginTop: 80 }}>暂无个人信息</p>

    return (
        <div
            style={{
                minHeight: '100vh',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                background: 'linear-gradient(135deg, #f5f7ff 0%, #ffffff 100%)',
                padding: '24px',
            }}
        >
            <Card
                style={{ width: 720, maxWidth: '100%', boxShadow: '0 8px 24px rgba(15, 23, 42, 0.08)' }}
                bodyStyle={{ padding: '24px 32px 28px' }}
            >
                <Space align="center" style={{ width: '100%', marginBottom: 24 }}>
                    <Avatar
                        size={72}
                        icon={<UserOutlined />}
                        style={{
                            background: 'linear-gradient(135deg, #6366f1 0%, #4f46e5 100%)',
                        }}
                    />
                    <div style={{ flex: 1 }}>
                        <Typography.Title level={3} style={{ margin: 0 }}>
                            {data.Nickname || '未设置昵称'}
                        </Typography.Title>
                        <Typography.Text type="secondary">
                            欢迎回来，这里是你的个人信息概览
                        </Typography.Text>
                    </div>
                    <Button href={"/users/edit"} type="primary">
                        编辑资料
                    </Button>
                </Space>

                <ProDescriptions
                    column={1}
                    size="middle"
                    bordered
                    colon={true}
                    labelStyle={{ width: 96 }}
                >
                    <ProDescriptions.Item label="邮箱" valueType="text">
                        <Space>
                            <MailOutlined />
                            {data.Email || <Typography.Text type="secondary">未绑定邮箱</Typography.Text>}
                        </Space>
                    </ProDescriptions.Item>
                    <ProDescriptions.Item label="手机" valueType="text">
                        <Space>
                            <PhoneOutlined />
                            {data.Phone || <Typography.Text type="secondary">未绑定手机</Typography.Text>}
                        </Space>
                    </ProDescriptions.Item>
                    <ProDescriptions.Item label="生日" valueType="text">
                        {data.Birthday || <Typography.Text type="secondary">未填写</Typography.Text>}
                    </ProDescriptions.Item>
                    <ProDescriptions.Item label="关于我" valueType="text">
                        {data.AboutMe
                            ? <Typography.Paragraph style={{ marginBottom: 0 }}>{data.AboutMe}</Typography.Paragraph>
                            : <Typography.Text type="secondary">还没有填写个人简介</Typography.Text>
                        }
                    </ProDescriptions.Item>
                </ProDescriptions>
            </Card>
        </div>
    )
}

export default Page