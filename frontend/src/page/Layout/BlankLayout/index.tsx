/*
 * Copyright (c) 2024 OceanBase.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React from 'react';
import { Layout } from '@oceanbase/design';
import useStyles from './index.style';

const BlankLayout: React.FC = ({ children, ...restProps }) => {
  const { styles } = useStyles();
  return (
    <div className={styles.main} {...restProps}>
      <Layout className={styles.layout}>{children}</Layout>
    </div>
  );
};

export default BlankLayout;
