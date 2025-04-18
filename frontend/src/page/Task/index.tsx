/*
 * Copyright (c) 2023 OceanBase
 * OCP Express is licensed under Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *          http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

import { formatMessage } from '@/util/intl';
import { history } from 'umi';
import React, { useState } from 'react';
import { Badge, Card, Radio, Space, Table, theme } from '@oceanbase/design';
import { some } from 'lodash';
import moment from 'moment';
import { PageContainer } from '@oceanbase/ui';
import { directTo, findByValue, sortByMoment } from '@oceanbase/util';
import { useUpdateEffect, useInterval, useRequest } from 'ahooks';
import { TASK_STATUS_LIST } from '@/constant/task';
import useDocumentTitle from '@/hook/useDocumentTitle';
import { formatTime } from '@/util/datetime';
import { getTaskProgress } from '@/util/task';
import MyInput from '@/component/MyInput';
import ContentWithReload from '@/component/ContentWithReload';
import { getAllAgentDags, getAllClusterDags } from '@/service/obshell/task';

export interface TaskProps {
  location?: {
    pathname: string;
    query: any;
  };

  // 展示模式: 页面模式 | 组件模式
  mode?: 'page' | 'component';
  // 任务类型: 手动触发的任务 | 定时调度的任务
  type?: API.TaskType;
  // 任务名称
  name?: string;
  clusterId?: number;
  tenantId?: number;
  hostId?: number;
}

const Task: React.FC<TaskProps> = ({
  location: { pathname, query = {} } = {},
  mode = 'page',
  type,
  name,
}) => {
  useDocumentTitle(formatMessage({ id: 'ocp-express.page.Task.Task', defaultMessage: '任务' }));
  const { token } = theme.useToken();

  const [status, setStatus] = useState(query.status || '');
  const [taskName, setTaskName] = useState(query.name);

  const { data, loading, refresh } = useRequest(getAllClusterDags, {});
  const {
    data: agentDagsRes,
    loading: agentDagsLoading,
    refresh: refreshAgentDags,
  } = useRequest(getAllAgentDags, {});

  const clusterTasks = data?.data?.contents || [];
  const agentTasks = agentDagsRes?.data?.contents || [];

  const realLoading = loading || agentDagsLoading;

  const sortTaskList = [...clusterTasks, ...agentTasks].sort((a, b) =>
    // 最近开始的任务排在前面
    sortByMoment(b, a, 'start_time')
  );

  const dataSource = sortTaskList.filter(item => {
    return (
      (!status || item.state === status) &&
      (!taskName || item.name?.toLowerCase()?.includes(taskName.toLowerCase()))
    );
  });

  // 是否应该轮询任务列表
  const polling = some(
    sortTaskList,
    // 列表中如果有运行中的任务，则会发起轮询，以保证状态能实时更新
    item => item.state === 'RUNNING'
  );

  useInterval(
    () => {
      refresh();
      refreshAgentDags();
    },
    polling ? 1000 : null
  );

  // 忽略首次渲染，只在依赖更新时调用，否则 initialPage 不生效
  useUpdateEffect(() => {
    if (mode === 'page') {
      // 搜索条件发生变化后重置到第一页
      // 搜索条件发生变化后更新 query 参数
      history.push({
        pathname,
        query: {
          status,
        },
      });
    }
  }, [type, status, name]);

  const columns = [
    {
      title: formatMessage({ id: 'ocp-express.page.Task.TaskName', defaultMessage: '任务名' }),
      dataIndex: 'name',
      render: (text: string, record: API.TaskInstance) => (
        <a
          data-aspm-click="c304181.d308752"
          data-aspm-desc="任务列表-跳转任务详情"
          data-aspm-param={``}
          data-aspm-expo
          onClick={() => {
            // component 表示任务列表内嵌到其他页面
            if (mode === 'component') {
              directTo(`/task/${record.id}?backUrl=/task`);
            } else {
              history.push(`/task/${record.id}`);
            }
          }}
        >
          {text}
        </a>
      ),
    },

    {
      title: 'ID',
      dataIndex: 'id',
    },

    {
      title: formatMessage({ id: 'ocp-express.page.Task.State', defaultMessage: '状态' }),
      dataIndex: 'state',
      render: (text: API.State) => {
        const statusItem = findByValue(TASK_STATUS_LIST, text);
        return <Badge status={statusItem.badgeStatus} text={statusItem.label} />;
      },
    },

    {
      title: formatMessage({
        id: 'ocp-express.page.Task.ExecutionProgress',
        defaultMessage: '执行进度',
      }),

      dataIndex: 'nodes',
      render: (text, record: API.DagDetailDTO) => <span>{getTaskProgress(record)}</span>,
    },

    {
      title: formatMessage({ id: 'ocp-express.page.Task.StartTime', defaultMessage: '开始时间' }),
      dataIndex: 'start_time',
      render: (text: string) => <span>{formatTime(text)}</span>,
    },

    {
      title: formatMessage({ id: 'ocp-express.page.Task.EndTime', defaultMessage: '结束时间' }),
      dataIndex: 'end_time',
      render: (text: string) => <span>{formatTime(text)}</span>,
    },
  ];

  const latestTime = moment().format('HH:mm:ss');
  const extra = (
    <Space size={16}>
      <Radio.Group
        value={status}
        onChange={e => {
          setStatus(e.target.value);
        }}
      >
        <Radio.Button value="">
          {formatMessage({ id: 'ocp-express.page.Task.All', defaultMessage: '全部' })}
        </Radio.Button>
        {TASK_STATUS_LIST.map(item => (
          <Radio.Button key={item.value} value={item.value}>
            {item.label}
          </Radio.Button>
        ))}
      </Radio.Group>
      {mode === 'page' && (
        <MyInput.Search
          data-aspm-click="c304181.d308754"
          data-aspm-desc="任务列表-搜索任务"
          data-aspm-param={``}
          data-aspm-expo
          // 由于是触发搜索后才修改 name 的值，因此这里使用 defaultValue，而不是 value，否则无法实时感知到最新的输入值
          defaultValue={query.name}
          onSearch={(value: string) => {
            setTaskName(value);
          }}
          placeholder={formatMessage({
            id: 'ocp-express.page.Task.SearchTaskName',
            defaultMessage: '搜索任务名称',
          })}
          className="search-input"
          allowClear={true}
        />
      )}

      <span
        style={{
          color: token.colorTextTertiary,
        }}
      >
        {formatMessage(
          {
            id: 'ocp-express.page.Task.LastRefreshTimeLatesttime',
            defaultMessage: '最近刷新时间：{latestTime}',
          },

          { latestTime }
        )}
      </span>
    </Space>
  );

  const content = (
    <div>
      <Table
        dataSource={dataSource}
        columns={columns}
        rowKey={record => record.id}
        loading={realLoading && !polling}
      />
    </div>
  );

  return mode === 'component' ? (
    <div>
      <div style={{ marginBottom: 24, textAlign: 'right' }}>{extra}</div>
      {content}
    </div>
  ) : (
    <PageContainer
      ghost={true}
      header={{
        title: (
          <ContentWithReload
            spin={loading && !polling}
            content={formatMessage({
              id: 'ocp-express.page.Task.TaskCenter',
              defaultMessage: '任务中心',
            })}
            onClick={() => {
              refresh();
            }}
          />
        ),
      }}
    >
      <Card
        className="card-without-padding"
        bordered={false}
        title={
          mode === 'page' &&
          formatMessage({
            id: 'ocp-express.page.Task.TaskList',
            defaultMessage: '任务列表',
          })
        }
        extra={extra}
      >
        {content}
      </Card>
    </PageContainer>
  );
};

export default Task;
